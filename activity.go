package sample

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"github.com/rstinear/workflow/client"
	r "github.com/rstinear/workflow/proto"
	md "google.golang.org/grpc/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New) //activity.Register(&Activity{}, New) to create instances using factory method 'New'
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

//New optional factory method, should be used if one activity instance per configuration is desired
func New(ctx activity.InitContext) (activity.Activity, error) {

	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		return nil, err
	}

	fmt.Println("settings: ", settings)

	act := &Activity{settings: settings} //add aSetting to instance

	return act, nil
}

// Activity is an sample Activity that can be used as a base to create a custom activity
type Activity struct {
	settings *Settings
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	input := &Input{}
	err = ctx.GetInputObject(input)

	fmt.Printf("my input value: %v\n", ctx.GetInput("projectID"))
	fmt.Printf("gidday, my input object is: %v and my settings object is: %v\n", input, a.settings)

	rc := client.New(a.settings.Token, a.settings.URI, true)

	gctx := md.AppendToOutgoingContext(context.Background(), "r-cs", input.Callsign, "r-p", input.ProjectID, "grpc-timeout", "0")

	msg := &r.CommandRequest{
		Txid:       generateTransactionID(),
		RequireAck: true,
		DeadlineMs: uint64(time.Now().Add(time.Second*10).UnixNano() / int64(time.Millisecond)),
		Uri: &r.RocosUri{
			Scheme: r.RocosUriScheme_GRPC,
			Path: &r.RocosUriPath{
				Project:   input.ProjectID,
				Callsign:  input.Callsign,
				Subsystem: "",
				Component: "command",
				Topic:     "do-that-thing",
			},
		},
	}

	fmt.Printf("about to command robot with msg %v\n", msg)
	cmdRespStream, err := rc.CallCommand(gctx, msg)
	if err != nil {
		fmt.Printf("failed to send command: %v\n", err)
		return false, err
	}

	var result string
o:
	for {
		resp, err := cmdRespStream.Recv()
		if err != nil {
			fmt.Printf("received err on ack stream: %v\n", err)
			return false, err
		}

		fmt.Printf("ack: %v\n", resp)

		switch m := resp.Content.(type) {
		case *r.ResponseChannelMessage_Heartbeat:
			fmt.Println("received heartbeat")
			continue
		case *r.ResponseChannelMessage_CommandResponse:
			fmt.Println("received command response", m)
			switch m.CommandResponse.Result.Status {
			case r.ResultStatus_COMPLETE_SUCCESS:
				result = "cool"
				break o
			default:
				fmt.Printf("received non success result code: %v", m.CommandResponse.Result.Status)
			}

		}

	}

	output := &Output{AnOutput: result}
	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}

func generateTransactionID() string {
	u := make([]byte, 16)
	_, err := rand.Read(u)
	if err != nil {
		log.Fatalf("could not generate transactionID: %v", err)
	}

	return hex.EncodeToString(u)
}

func gtTimestampMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
