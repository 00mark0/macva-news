// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.857
// components/adminLayout.templ

package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/00mark0/macva-news/db/services"
import "time"
import "fmt"

func AdminLayout(payload db.GetUserByIDRow, children ...templ.Component) templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"><title>Mačva News - Admin Panel</title><meta name=\"description\" content=\"Mačva News Admin Panel\"><link rel=\"preload\" href=\"/static/css/output.css\" as=\"style\"><link rel=\"preload\" href=\"/static/js/htmx.min.js\" as=\"script\"><link rel=\"preload\" href=\"/static/js/dark-mode.js\" as=\"script\"><link rel=\"preload\" href=\"/static/react/main.js\" as=\"script\"><link rel=\"preload\" href=\"/static/js/admin.js\" as=\"script\"><link href=\"/static/css/output.css\" rel=\"stylesheet\"><script src=\"/static/js/htmx.min.js\" defer></script><script src=\"/static/js/dark-mode.js\" defer></script><script src=\"/static/js/admin.js\" defer></script><script src=\"/static/react/main.js\" defer></script></head><body><nav class=\"fixed top-0 z-50 w-full bg-white border-b border-gray-200 dark:bg-black dark:border-gray-200\"><div class=\"px-3 py-3 lg:px-5 lg:pl-3\"><div class=\"flex items-center justify-between\"><div class=\"flex items-center justify-start rtl:justify-end\"><button data-drawer-target=\"logo-sidebar\" data-drawer-toggle=\"logo-sidebar\" aria-controls=\"logo-sidebar\" type=\"button\" class=\"inline-flex items-center p-2 text-sm text-gray-500 rounded-lg sm:hidden hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-200 dark:text-gray-400 dark:hover:bg-gray-700 dark:focus:ring-gray-600\"><span class=\"sr-only\">Open sidebar</span> <svg class=\"w-6 h-6\" aria-hidden=\"true\" fill=\"currentColor\" viewBox=\"0 0 20 20\" xmlns=\"http://www.w3.org/2000/svg\"><path clip-rule=\"evenodd\" fill-rule=\"evenodd\" d=\"M2 4.75A.75.75 0 012.75 4h14.5a.75.75 0 010 1.5H2.75A.75.75 0 012 4.75zm0 10.5a.75.75 0 01.75-.75h7.5a.75.75 0 010 1.5h-7.5a.75.75 0 01-.75-.75zM2 10a.75.75 0 01.75-.75h14.5a.75.75 0 010 1.5H2.75A.75.75 0 012 10z\"></path></svg></button> <a href=\"/admin\" class=\"flex ms-2 md:me-24\"><img src=\"/static/assets/f72dd544-c70b-4e10-a211-6c07a9478b44-macva-news-logo-cropped.webp\" class=\"sm:h-14 sm:w-64 h-12 w-48 object-fit rounded-lg me-3\" alt=\"Macva News Logo\"><h1 class=\"sm:block hidden self-center text-xl font-semibold whitespace-nowrap dark:text-white\">Mačva News - Admin Panel</h1></a></div><div class=\"flex items-center\"><div class=\"flex items-center ms-3\"><div><button type=\"button\" class=\"cursor-pointer flex text-sm rounded-full focus:ring-4 focus:ring-gray-300 dark:focus:ring-gray-600\" aria-expanded=\"false\" data-dropdown-toggle=\"dropdown-user\"><span class=\"sr-only\">Open user menu</span> <img class=\"w-10 h-10 rounded-full object-fit\" src=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(payload.Pfp)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/adminLayout.templ`, Line: 77, Col: 28}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "\" alt=\"user photo\" onerror=\"this.onerror=null; this.src=&#39;/static/assets/default-avatar-64x64.png&#39;;\"></button></div><div class=\"z-50 hidden absolute top-12 right-0 w-48 text-base list-none bg-white divide-y divide-gray-100 rounded-sm shadow-sm dark:bg-black dark:divide-gray-200\" id=\"dropdown-user\"><div class=\"px-4 py-3\" role=\"none\"><p class=\"text-sm text-gray-900 dark:text-white\" role=\"none\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(payload.Username)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/adminLayout.templ`, Line: 89, Col: 29}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "</p><p class=\"text-sm font-medium text-gray-900 truncate dark:text-gray-300\" role=\"none\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(payload.Email)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/adminLayout.templ`, Line: 92, Col: 26}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "</p></div><ul class=\"py-1\" role=\"none\"><li class=\"cursor-pointer\"><a id=\"user-menu-item-overview\" hx-trigger=\"click\" hx-get=\"/admin/hx-admin\" hx-target=\"#admin-content\" hx-swap=\"innerHTML\" class=\"block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-600 dark:hover:text-white\" role=\"menuitem\">Analitika</a></li><li class=\"cursor-pointer\"><a hx-get=\"/admin/settings\" hx-trigger=\"click\" hx-target=\"#admin-content\" hx-swap=\"innerHTML\" id=\"user-menu-item-settings\" class=\"block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-600 dark:hover:text-white\" role=\"menuitem\">Podešavanja</a></li><li><a href=\"#\" id=\"user-menu-item-earnings\" class=\"block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-600 dark:hover:text-white\" role=\"menuitem\">Earnings</a></li><li><a href=\"/\" class=\"cursor-pointer block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-600 dark:hover:text-white\" role=\"menuitem\">Naslovna</a></li><li><a id=\"user-menu-item-logout\" hx-post=\"/api/logout\" class=\"cursor-pointer block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-600 dark:hover:text-white\" role=\"menuitem\">Odjavi se</a></li></ul></div></div></div></div></div></nav><aside id=\"logo-sidebar\" class=\"fixed top-0 left-0 z-40 w-64 h-screen pt-20 transition-transform -translate-x-full bg-white border-r border-gray-200 sm:translate-x-0 dark:bg-black dark:border-gray-200\" aria-label=\"Sidebar\"><div class=\"h-full px-3 pb-4 overflow-y-auto bg-white dark:bg-black\"><ul class=\"space-y-2 font-medium pt-3\"><li class=\"cursor-pointer\"><a id=\"pregled\" hx-trigger=\"click\" hx-get=\"/admin/hx-admin\" hx-target=\"#admin-content\" hx-swap=\"innerHTML\" class=\"flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group\"><svg class=\"w-5 h-5 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white\" aria-hidden=\"true\" xmlns=\"http://www.w3.org/2000/svg\" fill=\"currentColor\" viewBox=\"0 0 22 21\"><path d=\"M16.975 11H10V4.025a1 1 0 0 0-1.066-.998 8.5 8.5 0 1 0 9.039 9.039.999.999 0 0 0-1-1.066h.002Z\"></path> <path d=\"M12.5 0c-.157 0-.311.01-.565.027A1 1 0 0 0 11 1.02V10h8.975a1 1 0 0 0 1-.935c.013-.188.028-.374.028-.565A8.51 8.51 0 0 0 12.5 0Z\"></path></svg> <span class=\"ms-3\">Analitika</span></a></li><li class=\"cursor-pointer\"><a id=\"kategorije\" hx-trigger=\"click\" hx-get=\"/admin/categories\" hx-target=\"#admin-content\" hx-swap=\"innerHTML\" class=\"flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group\"><svg class=\"shrink-0 w-5 h-5 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white\" aria-hidden=\"true\" xmlns=\"http://www.w3.org/2000/svg\" fill=\"currentColor\" viewBox=\"0 0 20 20\"><path d=\"M5 5V.13a2.96 2.96 0 0 0-1.293.749L.879 3.707A2.96 2.96 0 0 0 .13 5H5Z\"></path> <path d=\"M6.737 11.061a2.961 2.961 0 0 1 .81-1.515l6.117-6.116A4.839 4.839 0 0 1 16 2.141V2a1.97 1.97 0 0 0-1.933-2H7v5a2 2 0 0 1-2 2H0v11a1.969 1.969 0 0 0 1.933 2h12.134A1.97 1.97 0 0 0 16 18v-3.093l-1.546 1.546c-.413.413-.94.695-1.513.81l-3.4.679a2.947 2.947 0 0 1-1.85-.227 2.96 2.96 0 0 1-1.635-3.257l.681-3.397Z\"></path> <path d=\"M8.961 16a.93.93 0 0 0 .189-.019l3.4-.679a.961.961 0 0 0 .49-.263l6.118-6.117a2.884 2.884 0 0 0-4.079-4.078l-6.117 6.117a.96.96 0 0 0-.263.491l-.679 3.4A.961.961 0 0 0 8.961 16Zm7.477-9.8a.958.958 0 0 1 .68-.281.961.961 0 0 1 .682 1.644l-.315.315-1.36-1.36.313-.318Zm-5.911 5.911 4.236-4.236 1.359 1.359-4.236 4.237-1.7.339.341-1.699Z\"></path></svg> <span class=\"flex-1 ms-3 whitespace-nowrap\">Kategorije</span></a></li><li class=\"cursor-pointer\"><a id=\"artikli\" hx-trigger=\"click\" hx-get=\"/admin/content\" hx-target=\"#admin-content\" hx-swap=\"innerHTML\" class=\"flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group\"><svg class=\"shrink-0 w-5 h-5 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white\" aria-hidden=\"true\" xmlns=\"http://www.w3.org/2000/svg\" fill=\"currentColor\" viewBox=\"0 0 20 20\"><path d=\"M5 5V.13a2.96 2.96 0 0 0-1.293.749L.879 3.707A2.96 2.96 0 0 0 .13 5H5Z\"></path> <path d=\"M6.737 11.061a2.961 2.961 0 0 1 .81-1.515l6.117-6.116A4.839 4.839 0 0 1 16 2.141V2a1.97 1.97 0 0 0-1.933-2H7v5a2 2 0 0 1-2 2H0v11a1.969 1.969 0 0 0 1.933 2h12.134A1.97 1.97 0 0 0 16 18v-3.093l-1.546 1.546c-.413.413-.94.695-1.513.81l-3.4.679a2.947 2.947 0 0 1-1.85-.227 2.96 2.96 0 0 1-1.635-3.257l.681-3.397Z\"></path> <path d=\"M8.961 16a.93.93 0 0 0 .189-.019l3.4-.679a.961.961 0 0 0 .49-.263l6.118-6.117a2.884 2.884 0 0 0-4.079-4.078l-6.117 6.117a.96.96 0 0 0-.263.491l-.679 3.4A.961.961 0 0 0 8.961 16Zm7.477-9.8a.958.958 0 0 1 .68-.281.961.961 0 0 1 .682 1.644l-.315.315-1.36-1.36.313-.318Zm-5.911 5.911 4.236-4.236 1.359 1.359-4.236 4.237-1.7.339.341-1.699Z\"></path></svg> <span class=\"flex-1 ms-3 whitespace-nowrap\">Artikli</span></a></li><li class=\"cursor-pointer\"><a id=\"korisnici\" hx-trigger=\"click\" hx-get=\"/admin/users\" hx-target=\"#admin-content\" hx-swap=\"innerHTML\" class=\"flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group\"><svg class=\"shrink-0 w-5 h-5 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white\" aria-hidden=\"true\" xmlns=\"http://www.w3.org/2000/svg\" fill=\"currentColor\" viewBox=\"0 0 20 18\"><path d=\"M14 2a3.963 3.963 0 0 0-1.4.267 6.439 6.439 0 0 1-1.331 6.638A4 4 0 1 0 14 2Zm1 9h-1.264A6.957 6.957 0 0 1 15 15v2a2.97 2.97 0 0 1-.184 1H19a1 1 0 0 0 1-1v-1a5.006 5.006 0 0 0-5-5ZM6.5 9a4.5 4.5 0 1 0 0-9 4.5 4.5 0 0 0 0 9ZM8 10H5a5.006 5.006 0 0 0-5 5v2a1 1 0 0 0 1 1h11a1 1 0 0 0 1-1v-2a5.006 5.006 0 0 0-5-5Z\"></path></svg> <span class=\"flex-1 ms-3 whitespace-nowrap\">Korisnici</span></a></li><li class=\"cursor-pointer\"><a id=\"reklame\" hx-trigger=\"click\" hx-get=\"/admin/ads\" hx-target=\"#admin-content\" hx-swap=\"innerHTML\" class=\"flex items-center p-2 text-gray-900 rounded-lg dark:text-white hover:bg-gray-100 dark:hover:bg-gray-700 group\"><svg class=\"shrink-0 w-5 h-5 text-gray-500 transition duration-75 dark:text-gray-400 group-hover:text-gray-900 dark:group-hover:text-white\" aria-hidden=\"true\" xmlns=\"http://www.w3.org/2000/svg\" fill=\"currentColor\" viewBox=\"0 0 18 20\"><path d=\"M17 5.923A1 1 0 0 0 16 5h-3V4a4 4 0 1 0-8 0v1H2a1 1 0 0 0-1 .923L.086 17.846A2 2 0 0 0 2.08 20h13.84a2 2 0 0 0 1.994-2.153L17 5.923ZM7 9a1 1 0 0 1-2 0V7h2v2Zm0-5a2 2 0 1 1 4 0v1H7V4Zm6 5a1 1 0 1 1-2 0V7h2v2Z\"></path></svg> <span class=\"flex-1 ms-3 whitespace-nowrap\">Oglasi</span></a></li></ul></div></aside><div id=\"admin-content\" class=\"sm:pl-64 pt-24 dark:bg-black\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, child := range children {
			templ_7745c5c3_Err = child.Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "</div><footer class=\"w-full bg-white p-2 dark:bg-black dark:text-gray-400\"><p class=\"block text-sm text-gray-500 text-center \">© ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var5 string
		templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(time.Now().Year()))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/adminLayout.templ`, Line: 291, Col: 39}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, " <a href=\"/\" class=\"hover:underline\">Mačva News™</a>. All Rights Reserved.</p></footer></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
