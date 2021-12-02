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
	First  Actuator
	Second Actuator
}

func (p *ActuatorPipeline) First(actuator Actuator) Pipeline {
	p.Actuator = actuator
	return p
}

func (p *ActuatorPipeline) Then(actuator Actuator) Pipeline {
	state := &ComposingActuator{
		First:  p.Actuator,
		Second: actuator,
	}
	return &ActuatorPipeline{
		Actuator: state,
	}

}

func (p *ActuatorPipeline) Exec(c *context.Context) error {
	return p.Actuator.Exec(c)
}

func (ps *ComposingActuator) Exec(c *context.Context) error {
	if err := ps.First.Exec(c); err != nil {
		return err
	}
	if err := ps.Second.Exec(c); err != nil {
		return err
	}
	return nil
}
