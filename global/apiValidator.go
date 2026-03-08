package global

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/es"
	"github.com/go-playground/locales/fr"
	"github.com/go-playground/locales/ko"
	"github.com/go-playground/locales/ru"
	"github.com/go-playground/locales/zh_Hans_CN"
	"github.com/go-playground/locales/zh_Hant"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	es_translations "github.com/go-playground/validator/v10/translations/es"
	fr_translations "github.com/go-playground/validator/v10/translations/fr"
	ko_translations "github.com/go-playground/validator/v10/translations/ko"
	ru_translations "github.com/go-playground/validator/v10/translations/ru"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	zh_tw_translations "github.com/go-playground/validator/v10/translations/zh_tw"
	"reflect"
)

func ApiInitValidator() {
	validate := validator.New(validator.WithRequiredStructEnabled())

	enT := en.New()
	cn := zh_Hans_CN.New()
	koT := ko.New()
	ruT := ru.New()
	esT := es.New()
	frT := fr.New()
	zhTwT := zh_Hant.New()

	uni := ut.New(enT, cn, koT, ruT, esT, frT, zhTwT)

	enTrans, _ := uni.GetTranslator("en")
	zhTrans, _ := uni.GetTranslator("zh_Hans_CN")
	koTrans, _ := uni.GetTranslator("ko")
	ruTrans, _ := uni.GetTranslator("ru")
	esTrans, _ := uni.GetTranslator("es")
	frTrans, _ := uni.GetTranslator("fr")
	zhTwTrans, _ := uni.GetTranslator("zh_Hant")

	trans_items := []struct {
		register func(*validator.Validate, ut.Translator) error
		trans    ut.Translator
	}{
		{zh_translations.RegisterDefaultTranslations, zhTrans},
		{en_translations.RegisterDefaultTranslations, enTrans},
		{ko_translations.RegisterDefaultTranslations, koTrans},
		{ru_translations.RegisterDefaultTranslations, ruTrans},
		{es_translations.RegisterDefaultTranslations, esTrans},
		{fr_translations.RegisterDefaultTranslations, frTrans},
		{zh_tw_translations.RegisterDefaultTranslations, zhTwTrans},
	}
	
	for _, trans_item := range trans_items {
		if err := trans_item.register(validate, trans_item.trans); err != nil {
			log.Fatalf("validator translation error: %v", err)
		}
	}
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			return field.Name
		}

		return T(label)
	})

	Validator.Validate = validate
	Validator.UT = uni 
	Validator.VTrans = zhTrans

	Validator.ValidStruct = func(ctx *gin.Context, i interface{}) []string {
		err := Validator.Validate.Struct(i)
		lang := ctx.GetHeader("Accept-Language")
		if lang == "" {
			lang = Config.Lang
		}
		trans := getTranslatorForLang(lang)
		return translateErrors(err, trans)
	}
	Validator.ValidVar = func(ctx *gin.Context, field interface{}, tag string) []string {
		err := Validator.Validate.Var(field, tag)
		lang := ctx.GetHeader("Accept-Language")
		if lang == "" {
			lang = Config.Lang
		}
		trans := getTranslatorForLang(lang)
		return translateErrors(err, trans)
	}
}

func translateErrors(err error, trans ut.Translator) []string {

	errList := make([]string, 0, 10)

	if err != nil {

		if _, ok := err.(*validator.InvalidValidationError); ok {
			errList = append(errList, err.Error())
			return errList
		}

		for _, e := range err.(validator.ValidationErrors) {
			errList = append(errList, e.Translate(trans))
		}
	}

	return errList
}


func getTranslatorForLang(lang string) ut.Translator {
	switch lang {
	case "zh_CN":
		fallthrough
	case "zh-CN":
		fallthrough
	case "zh":
		trans, _ := Validator.UT.GetTranslator("zh_Hans_CN")
		return trans
	case "zh_TW":
		fallthrough
	case "zh-TW":
		fallthrough
	case "zh-tw":
		trans, _ := Validator.UT.GetTranslator("zh_Hant")
		return trans
	case "ko":
		trans, _ := Validator.UT.GetTranslator("ko")
		return trans
	case "ru":
		trans, _ := Validator.UT.GetTranslator("ru")
		return trans
	case "es":
		trans, _ := Validator.UT.GetTranslator("es")
		return trans
	case "fr":
		trans, _ := Validator.UT.GetTranslator("fr")
		return trans
	case "en":
		fallthrough
	default:
		trans, _ := Validator.UT.GetTranslator("en")
		return trans
	}
}
