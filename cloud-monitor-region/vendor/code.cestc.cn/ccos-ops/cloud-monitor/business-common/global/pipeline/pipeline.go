package pipeline

import "context"

type Actuator func(c *context.Context) error

type Pipeline interface {
	First(Actuator) Pipeline
	Then(Actuator) Pipeline
	Exec(*context.Context) error
}

type ActuatorPipeline struct {
	Actuator Actuator
}

func (p *ActuatorPipeline) First(ac Actuator) Pipeline {
	p.Actuator = ac
	return p
}

func (p *ActuatorPipeline) Then(ac Actuator) Pipeline {
	return &ActuatorPipeline{
		Actuator: func(c *context.Context) error {
			if err := p.Actuator(c); err != nil {
				return err
			}
			if err := ac(c); err != nil {
				return err
			}
			return nil
		},
	}

}

func (p *ActuatorPipeline) Exec(c *context.Context) error {
	return p.Actuator(c)
}
