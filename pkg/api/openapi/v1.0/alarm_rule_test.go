package v1_0

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"testing"
)

type User struct {
	FirstName      string     `binding:"oneof=enabled disabled"`
	LastName       string     `validate:"required"`
	Age            uint8      `validate:"gte=0,lte=130"`
	Email          string     `validate:"required,email"`
	FavouriteColor string     `validate:"iscolor"`                // alias for 'hexcolor|rgb|rgba|hsl|hsla'
	Addresses      []*Address `validate:"required,dive,required"` // a person can have a home and cottage...
}

// Address houses a users address information
type Address struct {
	Street string `validate:"required"`
	City   string `validate:"required"`
	Planet string `validate:"required"`
	Phone  string `validate:"required"`
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func TestVale1(t *testing.T) {

	validate = validator.New()

	//validateStruct()
	validateVariable()
}

func validateStruct() {

	address := &Address{
		Street: "Eavesdown Docks",
		Planet: "Persphone",
		Phone:  "none",
	}

	user := &User{
		FirstName:      "Badger",
		LastName:       "Smith",
		Age:            135,
		Email:          "Badger.Smith@gmail.com",
		FavouriteColor: "#000-",
		Addresses:      []*Address{address},
	}

	// returns nil or ValidationErrors ( []FieldError )
	err := validate.Struct(user)
	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}

		for _, err := range err.(validator.ValidationErrors) {

			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}

		// from here you can create your own error messages in whatever language you wish
		return
	}

	// save user to database
}

func TestVale2(t *testing.T) {

	validate = validator.New()
	validateVariable()
}
func validateVariable() {

	myEmail := ""

	errs := validate.Var(myEmail, "oneof=email d")

	if errs != nil {
		fmt.Println(errs) // output: Key: "" Error:Field validation for "" failed on the "email" tag
		return
	}

	// email ok, move on
}

type Stu struct {
	FirstName string `validate:"oneof=enabled disabled"`
}

func TestVale3(t *testing.T) {
	validate = validator.New()
	user := &Stu{}
	err := validate.Struct(user)
	fmt.Println(err)
}

func TestVale13(t *testing.T) {
	print(fmt.Sprintf("%.2f", 12345.777))
}

func TestVale14(t *testing.T) {
	var q []*User
	var x interface{}
	x = q
	user, ok := x.([]*Stu)
	if ok {

	}
	fmt.Println(user)
	for _, u := range user {
		fmt.Println(u)
	}

}

func TestVale5(t *testing.T) {
	reqParam := AlarmRuleCreateReqDTO{
		MonitorType:         "云产品监控",
		ProductAbbreviation: "ecs",
		Scope:               "INSTANCE",
		Resources: []struct {
			ResourceId string `binding:"required"`
		}{
			{ResourceId: "11"}, {ResourceId: "22"},
		},
		RuleName: "TEST1",
		TriggerCondition: struct {
			MetricCode         string `binding:"required"`
			Period             int    `binding:"required"`
			Times              int    `binding:"required"`
			Statistics         string `binding:"oneof=Maximum Minimum Average"`
			ComparisonOperator string `binding:"oneof=greater greaterOrEqual less  lessOrEqual  equal notEqual"`
			Threshold          string
		}{
			MetricCode:         "ecs_cpu_usage",
			Period:             60,
			Times:              1,
			Statistics:         "Maximum",
			ComparisonOperator: "greaterOrEqual",
			Threshold:          "10",
		},
		SilencesTime: "",
		AlarmLevel:   1,
		GroupList:    []string{"1", "2"},
	}
	marshal, _ := json.Marshal(reqParam)
	fmt.Println(string(marshal))
}

func TestVale15(t *testing.T) {

}
