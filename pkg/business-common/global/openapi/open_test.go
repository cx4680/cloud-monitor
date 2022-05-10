package openapi

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
	"testing"
	"time"
)

// DbBackedUser User struct
type DbBackedUser struct {
	Name sql.NullString `validate:"required"`
	Age  sql.NullInt64  `validate:"required"`
}

// use a single instance of Validate, it caches struct info

func TestField(t *testing.T) {

	validate := validator.New()

	// register all sql.Null* types to use the ValidateValuer CustomTypeFunc
	validate.RegisterCustomTypeFunc(ValidateValuer, sql.NullString{}, sql.NullInt64{}, sql.NullBool{}, sql.NullFloat64{})

	// build object for validation
	x := DbBackedUser{Name: sql.NullString{String: "www", Valid: true}, Age: sql.NullInt64{Int64: 3, Valid: true}}

	err := validate.Struct(x)

	if err != nil {
		fmt.Printf("Err(s):\n%+v\n", err)
	}
}

// ValidateValuer implements validator.CustomTypeFunc
func ValidateValuer(field reflect.Value) interface{} {

	if valuer, ok := field.Interface().(driver.Valuer); ok {

		val, err := valuer.Value()
		if err == nil {
			return val
		}
		// handle the error how you want
	}

	return nil
}

type MyStruct struct {
	String string `validate:"is-awesome"`
}

func TestField3(t *testing.T) {

	validate := validator.New()
	validate.RegisterValidation("is-awesome", ValidateMyVal)

	s := MyStruct{String: "awesome"}

	err := validate.Struct(s)
	if err != nil {
		fmt.Printf("Err(s):\n%+v\n", err)
	}

	s.String = "not awesome"
	err = validate.Struct(s)
	if err != nil {
		fmt.Printf("Err(s):\n%+v\n", err)
	}
}

// ValidateMyVal implements validator.Func
func ValidateMyVal(fl validator.FieldLevel) bool {
	return fl.Field().String() == "awesome"
}

var (
	Uni      *ut.UniversalTranslator
	Validate *validator.Validate
)

type User struct {
	Username  string    `form:"user_name" binding:"required"`
	Tagline   string    `form:"tag_line" binding:"required,lt=10"`
	Tagline2  string    `form:"tag_line2" binding:"required,oneof=1 2 3"`
	StartTime time.Time `time_format:"2006-01-02 15:04:05"`
	EndTime   time.Time `binding:"gtefield=StartTime" time_format:"2006-01-02 15:04:05"`
}

func TestField4(t *testing.T) {
	en := en.New()
	zh := zh.New()
	zh_tw := zh_Hant_TW.New()
	Uni = ut.New(en, zh, zh_tw)
	Validate = binding.Validator.Engine().(*validator.Validate)

	route := gin.Default()
	route.GET("/5lmh", startPage)
	route.POST("/5lmh", startPage)
	route.Run(":8081")
}

func startPage(c *gin.Context) {
	//这部分应放到中间件中
	locale := c.DefaultQuery("locale", "zh")
	trans, _ := Uni.GetTranslator(locale)
	switch locale {
	case "zh":
		zhTranslations.RegisterDefaultTranslations(Validate, trans)
		break
	case "en":
		enTranslations.RegisterDefaultTranslations(Validate, trans)
		break
	default:
		zhTranslations.RegisterDefaultTranslations(Validate, trans)
		break
	}

	//自定义错误内容
	Validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "the request parameter  {0} is not supplied!", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	//自定义错误内容
	Validate.RegisterTranslation("oneof", trans, func(ut ut.Translator) error {
		return ut.Add("oneof", "The required parameter {0} is not valid!", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("oneof", fe.Field())
		return t
	})

	//这块应该放到公共验证方法中
	user := User{}
	err := c.ShouldBind(&user)
	if parseError, ok := err.(*time.ParseError); ok {
		c.String(200, "The required parameter {0} is not valid!", parseError)
		return
	}

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range errs {
			fieldError.StructField()
			fieldError.Param()
		}
	}

	if err != nil {
		errs := err.(validator.ValidationErrors)
		var sliceErrs []string
		for _, e := range errs {
			sliceErrs = append(sliceErrs, e.Translate(trans))
		}
		c.String(200, fmt.Sprintf("%#v", sliceErrs))
		return
	}
	c.String(200, fmt.Sprintf("%#v", "user"))
}

func TestUri(t *testing.T) {
	uri := "/a/v/c/d"
	fmt.Println(strings.HasPrefix(uri, "/a/v"))
}

type Person struct {
	ID        int
	FirstName string `json:"name"`
	LastName  string
	Address   string `json:"address,omitempty"`
}

type Employee struct {
	Person
	ManagerID int
}

type Contractor struct {
	Person
	CompanyID int
}

func TestUri2(t *testing.T) {
	employees := []Employee{
		Employee{
			Person: Person{
				LastName: "Doe", FirstName: "John",
			},
			ManagerID: 333,
		},
		Employee{
			Person: Person{
				LastName: "Campbell", FirstName: "David",
			},
			ManagerID: 22,
		},
	}

	data, _ := json.Marshal(employees)
	fmt.Printf("%s\n", data)

	var decoded []Employee
	json.Unmarshal(data, &decoded)
	fmt.Printf("%v", decoded)
}
