package validation

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/locales/ja"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	jaTranslations "github.com/go-playground/validator/v10/translations/ja"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type CustomValidator struct {
	trans     ut.Translator
	validator *validator.Validate
}

func InitValidator() echo.Validator {
	jaTrans := ja.New()
	uniTrans := ut.New(jaTrans, jaTrans)
	trans, _ := uniTrans.GetTranslator("ja")

	validate := validator.New()

	// エラーメッセージの日本語化
	if err := jaTranslations.RegisterDefaultTranslations(validate, trans); err != nil {
		log.Fatal(err)
	}

	// フィールド名の日本語化（jaタグを登録）
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		fieldName := field.Tag.Get("ja")
		if fieldName == "-" {
			return ""
		}
		return fieldName
	})

	// カスタム型を登録
	validate.RegisterCustomTypeFunc(ValidateUuidValuer, uuid.UUID{})
	validate.RegisterCustomTypeFunc(ValidateSqlValuer, sql.NullString{}, sql.NullInt64{}, sql.NullBool{}, sql.NullFloat64{})

	// カスタムバリデーションを登録
	if err := validate.RegisterValidation("is-messi", ValidateIsMessi); err != nil {
		log.Fatal(err)
	}
	// 郵便番号チェック
	if err := validate.RegisterValidation("postal_code", PostalCodeValidation); err != nil {
		log.Fatal(err)
	}
	// 比較対象がnilの時のless than
	if err := validate.RegisterValidation("ltfield_if_Max_is_explicit", ValidateLessThanIfMaxIsExplicit); err != nil {
		log.Fatal(err)
	}
	// 比較対象がnilの時のless than
	if err := validate.RegisterValidation("exculded_if_for_bool", ValidateExcludedIfForBool); err != nil {
		log.Fatal(err)
	}

	// 日付(yyyy-MM-dd形式)チェックエラーメッセージ
	if err := validate.RegisterTranslation(
		"postal_code",
		trans,
		registerTranslator("postal_code", "{0}は正しい郵便番号の形式で指定してください"),
		translate,
	); err != nil {
		log.Fatal(err)
	}

	// ltnilfieldチェックエラーメッセージ
	if err := validate.RegisterTranslation(
		"ltfield_if_Max_is_explicit",
		trans,
		registerTranslator("ltfield_if_Max_is_explicit", "{0}は最大数よりも小さくなければなりません"),
		translate,
	); err != nil {
		log.Fatal(err)
	}
	// ltnilfieldチェックエラーメッセージ
	if err := validate.RegisterTranslation(
		"exculded_if_for_bool",
		trans,
		registerTranslator("exculded_if_for_bool", "{0}は入力不要です"),
		translate,
	); err != nil {
		log.Fatal(err)
	}

	// カスタムバリデーションの日本語エラーメッセージを登録
	if err := validate.RegisterTranslation("is-messi", trans, func(ut ut.Translator) error {
		return ut.Add("is-messi", "{0}はmessi以外認めません", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("is-messi", fe.Field())

		return t
	}); err != nil {
		log.Fatal(err)
	}

	return &CustomValidator{
		trans:     trans,
		validator: validate,
	}
}

func PostalCodeValidation(fl validator.FieldLevel) bool {
	regex := regexp.MustCompile(`^\d{3}-?\d{4}$`)
	return regex.MatchString(fl.Field().String())
}

func registerTranslator(tag string, msg string) validator.RegisterTranslationsFunc {
	return func(trans ut.Translator) error {
		if err := trans.Add(tag, msg, false); err != nil {
			return err
		}
		return nil
	}
}

func translate(trans ut.Translator, fe validator.FieldError) string {
	msg, err := trans.T(fe.Tag(), fe.Field())
	if err != nil {
		panic(fe.(error).Error())
	}
	return msg
}

// カスタムバリデーション
func ValidateIsMessi(fl validator.FieldLevel) bool {
	return fl.Field().String() == "messi"
}

// カスタムバリデーション
func ValidateLessThanIfMaxIsExplicit(fl validator.FieldLevel) bool {
	curr := fl.Field()
	cmp := fl.Parent().FieldByName(fl.Param())
	cmpType := cmp.Kind().String()

	switch {
	case cmpType == "ptr" && cmp.IsNil():
		return true
	case cmpType == "ptr" && (curr.Int() < cmp.Elem().Int()):
		return true
	case cmp.CanInt() && (curr.Int() < cmp.Int()):
		return true
	default:
		return false
	}
}

// カスタムバリデーション
func ValidateExcludedIfForBool(fl validator.FieldLevel) bool {
	paramSlice := strings.Split(fl.Param(), " ")
	cmpField := fl.Parent().FieldByName(paramSlice[0])
	boolStr := paramSlice[1]

	if boolStr == strconv.FormatBool(cmpField.Elem().Bool()) {
		return fl.Field().IsZero()
	}

	return true
}

// uuid.UUID型を登録
func ValidateUuidValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(uuid.UUID); ok {
		val, err := valuer.Value()
		if err == nil {
			return val
		}
		// handle the error how you want
	}

	return nil
}

// ValidateValuer implements validator.CustomTypeFunc
func ValidateSqlValuer(field reflect.Value) interface{} {

	if valuer, ok := field.Interface().(driver.Valuer); ok {

		val, err := valuer.Value()
		if err == nil {
			return val
		}
		// handle the error how you want
	}

	return nil
}

// バリデーション処理
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		var messages string
		for _, m := range err.(validator.ValidationErrors).Translate(cv.trans) {
			messages += m + "; "
		}
		// Optionally, you could return the error to give each route more control over the status code
		return errors.New(strings.TrimSuffix(messages, "; "))
	}
	return nil
}
