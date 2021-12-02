package pipeline

import "context"

type Actuator interface {
	Exec(c *context.Context) error
}

type Pipeline interface {
	Actuator
	First(Actuator) Pipeline
	Then(Actuator) Pipeline
}

type ActuatorPipeline struct {
	Actuator Actuator
}

type ComposingActuator struct {
	F Actuator
	S Actuator
}

func (p *ActuatorPipeline) First(a Actuator) Pipeline {
	p.Actuator = a
	return p
}

func (p *ActuatorPipeline) Then(a Actuator) Pipeline {
	state := &ComposingActuator{
		F: p.Actuator,
		S: a,
	}
	return &ActuatorPipeline{
		Actuator: state,
	}

}

func (p *ActuatorPipeline) Exec(c *context.Context) error {
	return p.Actuator.Exec(c)
}

func (ps *ComposingActuator) Exec(c *context.Context) error {
	if err := ps.F.Exec(c); err != nil {
		return err
	}
	if err := ps.S.Exec(c); err != nil {
		return err
	}
	return nil
}
