package sample

import "github.com/project-flogo/core/data/coerce"

type Settings struct {
	Token string `md:"token,required"`
	URI   string `md:"uri,required"`
}

type Input struct {
	ProjectID string `md:"projectID,required"`
	Callsign  string `md:"callsign,required"`
}

func (r *Input) FromMap(values map[string]interface{}) error {
	projectID, _ := coerce.ToString(values["projectID"])
	r.ProjectID = projectID

	callsign, _ := coerce.ToString(values["callsign"])
	r.Callsign = callsign
	return nil
}

func (r *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"projectID": r.ProjectID,
		"callsign":  r.Callsign,
	}
}

type Output struct {
	AnOutput string `md:"anOutput"`
}

func (o *Output) FromMap(values map[string]interface{}) error {
	strVal, _ := coerce.ToString(values["anOutput"])
	o.AnOutput = strVal
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"anOutput": o.AnOutput,
	}
}
