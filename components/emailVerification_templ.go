// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.857
// verification_success.templ

package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func VerificationSuccess() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!doctype html><html lang=\"sr\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>Uspešna Verifikacija - Mačva News</title><link rel=\"preload\" href=\"/static/css/output.css\" as=\"style\"><link href=\"/static/css/output.css\" rel=\"stylesheet\"></head><body class=\"bg-gray-100 min-h-screen flex items-center justify-center\"><div class=\"bg-white rounded-lg shadow-xl p-8 max-w-md w-full text-center\"><div class=\"mb-6\"><div class=\"mx-auto w-24 h-24 bg-green-100 rounded-full flex items-center justify-center\"><svg xmlns=\"http://www.w3.org/2000/svg\" class=\"h-16 w-16 text-green-500\" fill=\"none\" viewBox=\"0 0 24 24\" stroke=\"currentColor\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M5 13l4 4L19 7\"></path></svg></div></div><h1 class=\"text-2xl font-bold text-gray-800 mb-4\">Email Uspešno Verifikovan!</h1><p class=\"text-gray-600 mb-8\">Vaša email adresa je uspešno verifikovana. Sada možete pristupiti svom nalogu sa svim funkcionalnostima.</p><a href=\"/login\" class=\"inline-block bg-blue-500 hover:bg-blue-600 text-white font-medium py-3 px-6 rounded-md transition duration-300\">Nazad na prijavu</a></div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func VerificationError() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "<!doctype html><html lang=\"sr\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>Greška pri Verifikaciji - Mačva News</title><link rel=\"preload\" href=\"/static/css/output.css\" as=\"style\"><link href=\"/static/css/output.css\" rel=\"stylesheet\"></head><body class=\"bg-gray-100 min-h-screen flex items-center justify-center\"><div class=\"bg-white rounded-lg shadow-xl p-8 max-w-md w-full text-center\"><div class=\"mb-6\"><div class=\"mx-auto w-24 h-24 bg-red-100 rounded-full flex items-center justify-center\"><svg xmlns=\"http://www.w3.org/2000/svg\" class=\"h-16 w-16 text-red-500\" fill=\"none\" viewBox=\"0 0 24 24\" stroke=\"currentColor\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M6 18L18 6M6 6l12 12\"></path></svg></div></div><h1 class=\"text-2xl font-bold text-gray-800 mb-4\">Greška pri Verifikaciji</h1><p class=\"text-gray-600 mb-8\">Nismo u mogućnosti da verifikujemo vaš email. Verifikacioni link je možda istekao ili je već iskorišćen. Pokušajte ponovo da se prijavite ili kontaktirajte našu podršku.</p><a href=\"/login\" class=\"inline-block bg-blue-500 hover:bg-blue-600 text-white font-medium py-3 px-6 rounded-md transition duration-300\">Nazad na prijavu</a></div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
