// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.857
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/00mark0/macva-news/db/services"
import "fmt"
import "github.com/00mark0/macva-news/utils"

func Search(content []db.SearchContentRow, searchCount int64, searchTerm string, nextLimit int, globalSettings db.GlobalSetting) templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div id=\"search-results\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = SearchResults(content, searchCount, searchTerm, nextLimit, globalSettings).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "</div><div id=\"other-content\" hx-get=\"/api/content/other\" hx-swap=\"innerHTML\" hx-trigger=\"load\" class=\"mt-16\"></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func SearchResults(content []db.SearchContentRow, searchCount int64, searchTerm string, nextLimit int, globalSettings db.GlobalSetting) templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "<div class=\"container mx-auto\"><div class=\"mb-6 flex flex-col sm:flex-row items-center justify-between\"><h1 class=\"text-xl font-bold text-gray-800 dark:text-gray-200\">Rezultati Pretrage za: \"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(searchTerm)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 17, Col: 102}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "\"</h1><p class=\"text-gray-600 dark:text-gray-300\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if searchCount == 0 {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "Nije pronađen nijedan rezultat")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else if nextLimit > int(searchCount) {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "Prikaz ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var4 string
			templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(searchCount))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 22, Col: 37}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, " od ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(searchCount))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 22, Col: 68}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, " rezultata")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 9, "Prikaz ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var6 string
			templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(nextLimit))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 24, Col: 35}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 10, " od ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var7 string
			templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(searchCount))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 24, Col: 66}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 11, " rezultata")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 12, "</p></div><div class=\"grid gap-6 md:grid-cols-2 lg:grid-cols-3\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, v := range content {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 13, "<article class=\"group relative bg-white shadow-md border-b-2 border-gray-200 transition-all duration-300 hover:shadow-lg hover:-translate-y-1 overflow-hidden dark:bg-gray-800 flex flex-col h-[364px]\" style=\"z-index: 10;\"><div class=\"flex flex-col h-full\"><!-- Image Container - Fixed height --><a href=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var8 templ.SafeURL = templ.SafeURL(fmt.Sprintf("/%s/%s", utils.Slugify(v.Title), v.ContentID.String()))
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var8)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 14, "\"><div class=\"relative w-full h-48\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if v.Thumbnail.Valid && v.Thumbnail.String != "" {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 15, "<img src=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var9 string
				templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs(v.Thumbnail.String)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 37, Col: 34}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 16, "\" alt=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var10 string
				templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(v.Title)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 38, Col: 23}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 17, "\" class=\"thumbnail absolute inset-0 w-full h-full object-cover transition-transform duration-500 group-hover:scale-105 brightness-90 group-hover:brightness-100\"><div class=\"absolute inset-0 bg-gradient-to-t from-black/80 via-black/30 to-transparent\"></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			} else {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 18, "<div class=\"absolute inset-0 bg-gradient-to-t from-gray-900 to-gray-700\"></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 19, "</div></a><!-- Content Container - Takes remaining height --><div class=\"px-4 pt-2 relative flex-grow flex flex-col justify-between z-10 bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm\"><div class=\"absolute top-0 left-0 w-1 h-full bg-gray-800 group-hover:bg-blue-600 transition-colors duration-300\"></div><!-- Top content section --><a href=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var11 templ.SafeURL = templ.SafeURL(fmt.Sprintf("/%s/%s", utils.Slugify(v.Title), v.ContentID.String()))
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var11)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 20, "\"><div><div class=\"flex justify-between items-center text-xs text-gray-600 dark:text-gray-300 mb-2 space-x-2\"><span class=\"dark:text-gray-300\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var12 string
			templ_7745c5c3_Var12, templ_7745c5c3_Err = templ.JoinStringErrs(utils.TimeAgo(v.PublishedAt.Time))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 54, Col: 78}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var12))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 21, "</span><p class=\"text-sm sm:text-md\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var13 string
			templ_7745c5c3_Var13, templ_7745c5c3_Err = templ.JoinStringErrs(v.CategoryName)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 55, Col: 56}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var13))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 22, "</p></div><h2 class=\"text-md font-semibold text-gray-900 dark:text-gray-200 mb-2 line-clamp-2 leading-tight\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var14 string
			templ_7745c5c3_Var14, templ_7745c5c3_Err = templ.JoinStringErrs(v.Title)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 58, Col: 19}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var14))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 23, "</h2><p class=\"text-sm text-gray-700 dark:text-gray-300 mb-2 line-clamp-2 leading-tight\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var15 string
			templ_7745c5c3_Var15, templ_7745c5c3_Err = templ.JoinStringErrs(utils.ParseHTMLToText(v.ContentDescription))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 61, Col: 55}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var15))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 24, "</p></div></a><!-- Bottom section with stats - Always at bottom --><div class=\"flex items-center justify-between mt-auto py-2 border-t border-gray-100 dark:border-gray-700\"><div class=\"flex space-x-3 text-xs text-gray-600 dark:text-gray-300\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if v.ViewCountEnabled && !globalSettings.DisableViews {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 25, "<span class=\"flex items-center hover:text-green-600 transition-colors\"><svg xmlns=\"http://www.w3.org/2000/svg\" class=\"h-4 w-4 mr-1\" viewBox=\"0 0 20 20\" fill=\"currentColor\"><path d=\"M10 12a2 2 0 100-4 2 2 0 000 4z\"></path> <path fill-rule=\"evenodd\" d=\"M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z\" clip-rule=\"evenodd\"></path></svg> ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var16 string
				templ_7745c5c3_Var16, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(v.ViewCount))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 74, Col: 36}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var16))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 26, "</span> ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if v.LikeCountEnabled && !globalSettings.DisableLikes {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 27, "<span class=\"flex items-center hover:text-red-600 transition-colors\"><svg xmlns=\"http://www.w3.org/2000/svg\" class=\"h-4 w-4 mr-1\" viewBox=\"0 0 20 20\" fill=\"currentColor\"><path d=\"M2 10.5a1.5 1.5 0 113 0v6a1.5 1.5 0 01-3 0v-6zM6 10.333v5.43a2 2 0 001.106 1.79l.05.025A4 4 0 008.943 18h5.416a2 2 0 001.962-1.608l1.2-6A2 2 0 0015.56 8H12V4a2 2 0 00-2-2 1 1 0 00-1 1v.667a4 4 0 01-.8 2.4L6.8 7.933a4 4 0 00-.8 2.4z\"></path></svg> ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var17 string
				templ_7745c5c3_Var17, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(v.LikeCount))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 82, Col: 36}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var17))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 28, "</span> ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if v.DislikeCountEnabled && !globalSettings.DisableDislikes {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 29, "<span class=\"flex items-center hover:text-red-600 transition-colors\"><svg xmlns=\"http://www.w3.org/2000/svg\" class=\"h-4 w-4 mr-1\" viewBox=\"0 0 20 20\" fill=\"currentColor\"><path d=\"M18 9.5a1.5 1.5 0 11-3 0v-6a1.5 1.5 0 013 0v6zM14 9.667v-5.43a2 2 0 00-1.105-1.79l-.05-.025A4 4 0 0011.055 2H5.64a2 2 0 00-1.962 1.608l-1.2 6A2 2 0 004.44 12H8v4a2 2 0 002 2 1 1 0 001-1v-.667a4 4 0 01.8-2.4l1.4-1.866a4 4 0 00.8-2.4z\"></path></svg> ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var18 string
				templ_7745c5c3_Var18, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(v.DislikeCount))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 90, Col: 39}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var18))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 30, "</span> ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if v.CommentsEnabled && !globalSettings.DisableComments {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 31, "<span class=\"flex items-center hover:text-blue-600 transition-colors\"><svg xmlns=\"http://www.w3.org/2000/svg\" class=\"h-4 w-4 mr-1\" viewBox=\"0 0 20 20\" fill=\"currentColor\"><path fill-rule=\"evenodd\" d=\"M18 10c0 3.866-3.582 7-8 7a8.841 8.841 0 01-4.083-.98L2 17l1.338-3.123C2.493 12.767 2 11.434 2 10c0-3.866 3.582-7 8-7s8 3.134 8 7zM7 9H5v2h2V9zm8 0h-2v2h2V9zM9 9h2v2H9V9z\" clip-rule=\"evenodd\"></path></svg> ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var19 string
				templ_7745c5c3_Var19, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(v.CommentCount))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 98, Col: 39}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var19))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 32, "</span>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 33, "</div></div></div></div></article>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 34, "</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if nextLimit < int(searchCount) {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 35, "<div class=\"text-center mt-8\"><button hx-get=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var20 string
			templ_7745c5c3_Var20, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("/api/search?search_term=%s&limit=%d", searchTerm, nextLimit))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 111, Col: 87}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var20))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 36, "\" hx-target=\"#search-results\" hx-swap=\"innerHTML\" hx-trigger=\"click\" class=\"cursor-pointer inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition duration-150 ease-in-out\">Učitaj više <svg xmlns=\"http://www.w3.org/2000/svg\" class=\"h-4 w-4 ml-2\" fill=\"none\" viewBox=\"0 0 24 24\" stroke=\"currentColor\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M19 9l-7 7-7-7\"></path></svg></button></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 37, "</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func OtherContent(content []ListPublishedContentRes, nextLimit int, globalSettings db.GlobalSetting) templ.Component {
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
		templ_7745c5c3_Var21 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var21 == nil {
			templ_7745c5c3_Var21 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 38, "<h3 class=\"text-xl font-bold mb-4 text-gray-800 dark:text-gray-200\">Nedavno</h3><div class=\"grid gap-6\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, v := range content {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 39, "<div class=\"group flex bg-white shadow-md border-b-2 border-gray-200 transition-all duration-300 hover:shadow-lg hover:-translate-y-1 dark:bg-gray-800 rounded-lg overflow-hidden\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if v.Thumbnail != "" {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 40, "<div class=\"w-1/4 relative h-full overflow-hidden\"><a href=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var22 templ.SafeURL = templ.SafeURL(fmt.Sprintf("/%s/%s", utils.Slugify(v.Title), v.ContentID))
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var22)))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 41, "\"><img src=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var23 string
				templ_7745c5c3_Var23, templ_7745c5c3_Err = templ.JoinStringErrs(v.Thumbnail)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 142, Col: 25}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var23))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 42, "\" alt=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var24 string
				templ_7745c5c3_Var24, templ_7745c5c3_Err = templ.JoinStringErrs(v.Title)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 143, Col: 21}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var24))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 43, "\" class=\"absolute inset-0 w-full h-full object-cover transition-transform duration-500 group-hover:scale-105 brightness-90 hover:brightness-100\"></a></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 44, "<div class=\"w-2/3 p-4 flex flex-col justify-between relative\"><div class=\"absolute top-0 left-0 w-1 h-full bg-gray-800 group-hover:bg-blue-600 transition-colors duration-300\"></div><a href=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var25 templ.SafeURL = templ.SafeURL(fmt.Sprintf("/%s/%s", utils.Slugify(v.Title), v.ContentID))
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var25)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 45, "\"><div><div class=\"w-full flex items-center text-xs text-gray-600 dark:text-gray-300 mb-2 space-x-2\"><span class=\"flex items-center\"><svg xmlns=\"http://www.w3.org/2000/svg\" class=\"h-3 w-3 mr-1 text-gray-400\" fill=\"none\" viewBox=\"0 0 24 24\" stroke=\"currentColor\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M5.121 17.804A13.937 13.937 0 0112 16c2.5 0 4.847.655 6.879 1.804M15 10a3 3 0 11-6 0 3 3 0 016 0zm6 2a9 9 0 11-18 0 9 9 0 0118 0z\"></path></svg> ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var26 string
			templ_7745c5c3_Var26, templ_7745c5c3_Err = templ.JoinStringErrs(v.Username)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 158, Col: 21}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var26))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 46, "</span> <span class=\"text-gray-300 dark:text-gray-600\">•</span> <span class=\"dark:text-gray-300 w-full\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var27 string
			templ_7745c5c3_Var27, templ_7745c5c3_Err = templ.JoinStringErrs(v.PublishedAt)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 161, Col: 63}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var27))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 47, "</span></div><h2 class=\"text-lg font-semibold text-gray-900 dark:text-gray-200 mb-2 line-clamp-2 leading-tight\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var28 string
			templ_7745c5c3_Var28, templ_7745c5c3_Err = templ.JoinStringErrs(v.Title)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 164, Col: 17}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var28))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 48, "</h2><p class=\"text-sm text-gray-700 dark:text-gray-300 mb-3 line-clamp-2 leading-tight\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var29 string
			templ_7745c5c3_Var29, templ_7745c5c3_Err = templ.JoinStringErrs(utils.ParseHTMLToText(v.ContentDescription))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 167, Col: 53}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var29))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 49, "</p></div></a><div class=\"flex items-center justify-between mt-2 pt-2 border-t border-gray-100\"><div class=\"flex space-x-3 text-xs text-gray-600 dark:text-gray-300\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if v.ViewCountEnabled && !globalSettings.DisableViews {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 50, "<span class=\"flex items-center hover:text-green-600 transition-colors\"><svg xmlns=\"http://www.w3.org/2000/svg\" class=\"h-4 w-4 mr-1\" viewBox=\"0 0 20 20\" fill=\"currentColor\"><path d=\"M10 12a2 2 0 100-4 2 2 0 000 4z\"></path> <path fill-rule=\"evenodd\" d=\"M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z\" clip-rule=\"evenodd\"></path></svg> ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var30 string
				templ_7745c5c3_Var30, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(v.ViewCount))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 179, Col: 34}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var30))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 51, "</span> ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if v.LikeCountEnabled && !globalSettings.DisableLikes {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 52, "<span class=\"flex items-center hover:text-red-600 transition-colors\"><svg xmlns=\"http://www.w3.org/2000/svg\" class=\"h-4 w-4 mr-1\" viewBox=\"0 0 20 20\" fill=\"currentColor\"><path d=\"M2 10.5a1.5 1.5 0 113 0v6a1.5 1.5 0 01-3 0v-6zM6 10.333v5.43a2 2 0 001.106 1.79l.05.025A4 4 0 008.943 18h5.416a2 2 0 001.962-1.608l1.2-6A2 2 0 0015.56 8H12V4a2 2 0 00-2-2 1 1 0 00-1 1v.667a4 4 0 01-.8 2.4L6.8 7.933a4 4 0 00-.8 2.4z\"></path></svg> ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var31 string
				templ_7745c5c3_Var31, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(v.LikeCount))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 187, Col: 34}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var31))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 53, "</span> ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if v.DislikeCountEnabled && !globalSettings.DisableDislikes {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 54, "<span class=\"flex items-center hover:text-red-600 transition-colors\"><svg xmlns=\"http://www.w3.org/2000/svg\" class=\"h-4 w-4 mr-1\" viewBox=\"0 0 20 20\" fill=\"currentColor\"><path d=\"M18 9.5a1.5 1.5 0 11-3 0v-6a1.5 1.5 0 013 0v6zM14 9.667v-5.43a2 2 0 00-1.105-1.79l-.05-.025A4 4 0 0011.055 2H5.64a2 2 0 00-1.962 1.608l-1.2 6A2 2 0 004.44 12H8v4a2 2 0 002 2 1 1 0 001-1v-.667a4 4 0 01.8-2.4l1.4-1.866a4 4 0 00.8-2.4z\"></path></svg> ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var32 string
				templ_7745c5c3_Var32, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(v.DislikeCount))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 195, Col: 37}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var32))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 55, "</span> ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if v.CommentsEnabled && !globalSettings.DisableComments {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 56, "<span class=\"flex items-center hover:text-blue-600 transition-colors\"><svg xmlns=\"http://www.w3.org/2000/svg\" class=\"h-4 w-4 mr-1\" viewBox=\"0 0 20 20\" fill=\"currentColor\"><path fill-rule=\"evenodd\" d=\"M18 10c0 3.866-3.582 7-8 7a8.841 8.841 0 01-4.083-.98L2 17l1.338-3.123C2.493 12.767 2 11.434 2 10c0-3.866 3.582-7 8-7s8 3.134 8 7zM7 9H5v2h2V9zm8 0h-2v2h2V9zM9 9h2v2H9V9z\" clip-rule=\"evenodd\"></path></svg> ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var33 string
				templ_7745c5c3_Var33, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(v.CommentCount))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 203, Col: 37}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var33))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 57, "</span>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 58, "</div></div></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 59, "</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if len(content) == nextLimit {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 60, "<div class=\"text-center mt-8\"><button hx-get=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var34 string
			templ_7745c5c3_Var34, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("/api/content/other?limit=%d", nextLimit))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/search.templ`, Line: 215, Col: 66}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var34))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 61, "\" hx-target=\"#other-content\" hx-swap=\"innerHTML\" hx-trigger=\"click\" class=\"cursor-pointer inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition duration-150 ease-in-out\">Učitaj više <svg xmlns=\"http://www.w3.org/2000/svg\" class=\"h-4 w-4 ml-2\" fill=\"none\" viewBox=\"0 0 24 24\" stroke=\"currentColor\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" stroke-width=\"2\" d=\"M19 9l-7 7-7-7\"></path></svg></button></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		return nil
	})
}

func SearchPage(props ...interface{}) templ.Component {
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
		templ_7745c5c3_Var35 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var35 == nil {
			templ_7745c5c3_Var35 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = Layout(props[0].(db.GetUserByIDRow), props[1].(Meta), props[2].([]db.Ad), props[3].([]db.Category), Search(props[4].([]db.SearchContentRow), props[5].(int64), props[6].(string), props[7].(int), props[8].(db.GlobalSetting))).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
