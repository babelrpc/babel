package babeltemplates

import (
	"github.com/golang/snappy"
	"encoding/base64"
	"log"
	"mime"
	"net/http"
	"path"
	"strings"
)

const (
	c_ErrorBabel = "9AhgLy8vIEJhYmVsJ3MgRXJyb3IgRm9ybWF0CgEZ4ENvcHlyaWdodCAoQykgMjAxMyBDb25jdXIgVGVjaG5vbG9naWVzLCBJbmMuCgpuYW1lc3BhY2UgYwUlEC5jb20vBVwECgoBTglfRGRlZmluZXMgYSBzaW5nbGUgZQV2WG1lc3NhZ2UgYW5kIGNvZGUgdGhhdCBtBX0sYmUgbG9jYWxpemVkBUoBJXBkaXNwbGF5ZWQgdG8gYSBjYWxsZXIuCnN0cnVjdA3LCHsKCQF5TFRoZSBzZXJ2aWNlLXNwZWNpZmljDXcBazgKCXN0cmluZyBDb2RlOwoZNCR0ZXh0IG9mIHRoEaYwaW4gVVMtRW5nbGlzaBU4AE0JvC47AAhsaXMFOyRwYXJhbWV0ZXJzAacZSQ3vWC4gVGhpcyBjb3VsZCBiZSB1c2VkIGJ5CYMN9CxhdGlvbiBzeXN0ZW0FRiBnZW5lcmF0ZSANRQxzIGJhATgEb24dqgHZCC4KCQGIADwJ3wg-IFABjQxzOwp9KZAAUykTNpcBGaMgcmVzcG9uc2UgEXABmkksDCBmb3IxUiAgZmFpbHVyZXMxezJXADGCCGltZTJKARgKCWRhdGV0ARcIVGltOXYEYWcF4yRjYXRlZ29yaXplGdM6zQAMVGFnczFyBEEgAegBYUVNAHNJPRxvY2N1cnJlZA05BegAPkkaFToIQ29uJeYAaUGRCG1hcAFEDGFkZGkhdARhbE0lHGRldGFpbHMgQXcAYwkxEAoJbWFwLVkELCAuDAAtbQA-YSoBWxGkAEQNRQFBGGFpbnMgb3ANZQlfBGVkDW4IaW5mZYIsaW9uLCBzdWNoIGFzDRsx8UkgBG9yYUUYdGFjayB0cmF3VdwEdG8lNxBnaXZlbiWoNGVyIGVudmlyb25tZW50VbkNlBGjKElubmVyIHBvaW50JYYAYQU8BGljUf5kcHJvcGFnYXRlZCBmcm9tIGFub3RoZXIgdGkBBWGwEGlzIHVzYUtp9wAJMgwCBV9Bfg=="
	cAsp_BaseAsp = "0AS8e3tkZWZpbmUgIkNPTU1FTlRTIn19e3skY210cyA6PSBleHBhbmRDb21tZW50cyAuAR0UcmFuZ2UgBSMBD1RpbmRlbnR9fScge3sufX0NCnt7ZW5kARoNBxFiDEFUVFINXwxhdHRyBWAcZmlsdGVyQXQBDwVdHGlmIGxlbiAkBSIBRRVfCFt7ew16EGksICR4AZ4dJhxmICRpfX0sIBV9EC5OYW1lBUEJUjguUGFyYW1ldGVyc319KHsRTBBqLCAkeQFMMiAAGHt7aWYgJGoyUQAIaWYgFVQNCQQgPRlwKGZvcm1hdFZhbHVlKScF_gApFZAFDwBdZh4BFE1FVEhPRJaGATGQAd8tFCGxAQ81L0KOAQkpPucALr0BKTIIIC0gNR4MaSwkbSUdLnIAJagFzz7OAAQNCg=="
	cAsp_ClientAsp = "4QzwSHt7aWYgaXNBc3B9fTwle3tlbmR9fXt7JGlkbCA6PSAufX0NCicgQVVUTy1HRU5FUkFURUQgRklMRSAtIERPIE5PVCBNT0RJRlkBJ2RHZW5lcmF0ZWQgZnJvbSB7ey5GaWxlbmFtZQFHVHt7dGVtcGxhdGUgIkNPTU1FTlRTIiABaiQuQ29tbWVudHMgASgBSqhSZXF1aXJlcyBpbmNfYmFiZWwuYXNwDQoNCnt7cmFuZ2UgLlNlcnZpY2VzBa8EZm4BrixmdWxsTmFtZU9mIC4BCAUbAHMFGxUQNHNldGluZGVudCAiXHQiASJOlgAALh2SSGNsYXNzIHt7JGZufX1DbGllbnQBiRAJUHJpdgHNEFVSTA0KFQ4oVGltZW91dFNlY3MdFhhIZWFkZXJzBTgQc3ViIEMBVFhfSW5pdGlhbGl6ZQ0KCQlVUkwgPSAiIgEMAFQZRhAgPSAzMAEUDHNldCANSEwgPSBOb3RoaW5nDQoJZW5kIHN1YjpcABxUZXJtaW5hdAVbnjsACFNldAliLChrZXksIHZhbHVlKQF8BGlmFXsEaXMRfBQgdGhlbiA2XACMQ3JlYXRlT2JqZWN0KCJTY3JpcHRpbmcuRGljdGlvbmFyeSIpAVEJZwBzAWgIKSA9CWo6zAAEJyAhIFANCgknIGJhc2VVcmwgLSBCYXNlIHMp-wAgKXoIJyB0NS4Yb25kcyAtIC0_BCBpAZ8AYwEVCA0KCSF4AVIQSHR0cCgNUgAsOjwABZwpiQ0iAa0NTyHUBCA9OjIAAA09gVGcFCRpLCAkbUGSJC5NZXRob2RzfX1Bum58AhRNRVRIT0R5GAwufX0JbZ10Vm9pZCAuUmV0dXJuc319c3Vie3tlbHNlfX1mdW5jIVdtuGFySfgAKDKLAAB2BYsgUGFyYW1ldGVyZTN1ETRpZiBsYXN0ICRpICRtLhklJCB8IG5vdH19LCARbg1mDQclKyhkaW0gcGFybXMgOkUJCQwEPSCWBwKCpQAAaWWqYaUdCjhmIC5UeXBlLklzU3RydWMBugxkbH19QYgNmyU4ESUITWFwMh0ABaMMKCJ7ey0CACJBgBUOECAnIHt7BTsps6UQAAlaigEMQ2FsbD0eGGlmIG9yICgxqjahAAApNcIILklzPpQAGX8APQFQEG5kfX1CgfUBYTgoVVJMLCAie3skc259fS8VKwQiLHGTACxR8EGhACwpkAUyLGludGVybmFsVHlwZRV_BH19ZYIMZW5kICUoBGlzomACQcQlEQQNCgE_pRmBdy53AQBpyUwEJT4NFQ=="
	cAsp_ModelAsp = "xR64e3tkZWZpbmUgIklOSVRGSUVMRFMifX17e2lmIGxlbiAuRmllbGRzfX0NCgkJJyAJDjAgZnJvbSB7ey5OYW1lAS0IZW5kAQcQcmFuZ2UZMwVFRC5UeXBlLklzU3RydWN0IGlkbAlMCHNldBlCLCA9IE5vdGhpbmd7ew1MCGlmIBE7CE1hcFY0AIxDcmVhdGVPYmplY3QoIlNjcmlwdGluZy5EaWN0aW9uYXJ5IilOUQAMTGlzdAlSFcQgID0gQXJyYXkoNjEAKEluaXRpYWxpemVyRjEAOHt7Zm9ybWF0VmFsdWUgLjIsABV3Dc8NBzFtEENMT1NF_m4B_m4Bum4BaqIBVlEBEEVtcHR5egkBPQQ6PgJADQp7e3NldGluZGVudCAiXHRFl3R0ZW1wbGF0ZSAiQ09NTUVOVFMiIC5Db21tZW50cyBBhRkiDEFUVFIBH0hBdHRyaWJ1dGVzfX0JcHVibGljXXsAJwEMQaUBPQxpZiAoUXcMRW51bUGwQCl9fSAtIHNlZSAie3tmdWxsQfYAT1GfAQ0gfX0iIGZvciB2QQkAcw3yAcAt9g0HMf0UVE9KU09OBc9G4QApBiAkaSwgJHYgOj15TgFQKQ0FqSUXSH19Y2FsbCBzXy5Xcml0ZShqXywBpRxpbnRlcm5hbAHWCCAkdgXRCH19IgUcBCR2bbYQIiwge3sVDgUcGHJlbmFtZXMyMwAEaV95EjrIAAxGUk9N_soAGcqFLgxvciAoDbEyNAQAKRGUCC5Jc4UPiT8NkRXGBCA9IQwMUmVhZP4LAUYLAUIHASxpZiBpc0FzcH19PCUt9BR7eyRpZGwlt4EvxCcgQVVUTy1HRU5FUkFURUQgRklMRSAtIERPIE5PVCBNT0RJRlkNCicgR2VuZXJhdGVktVoIRmlsJY0BRwR7e1nUAENV9gFqLvoCBA0KRSKtt0G9IHN9fScgKioqIAENAWAMaW9ucwERCSstPUlxAC4BJABzQVoIJGZuAcEu6QJJ5DKPA0W-TpcALpMABHt7DV2lIhhzfX1jb25zofcIJGZuAWtNigQgPUG1DQ0MIiAnIAUvDCBvZiA6YgUpDQW4EA0KZnVupeZBxgFPjF9Ub0lkKHN0clZhbCkNCglkaW0gcg0KCXNlbGVjdCBjYXNlIAkeISg6mQAACQUfACLVDwAiwe8IciA9aoIACTMMZWxzZREsAGWB8BgNCgllbmQgCXYQDQoJe3sF7QWeBCA9AZIBHxG9IYoRDB3JIFN0cmluZyhpZGrMAAUdbssAMp0GBH19FdMy5gBNucLNAAmiTtEADVIBmh5ECABDIfdVjgUOAGGhrwgqKiohAQ0yVS4FIf6MAuaMAm51Ag2rLt0AFtkIFd4JDwh1cmViagMJIn7gABQkYmFzZXNhiiAuQmFzZUNsYXMBEGH2YTsMJHN1YgkeCFN1Yj4dAN1WmqgDecpONQcAYwF9UZ-lqA4bCgmiYZQ8CScgaW5oZXJpdGFuY2U6IDHKxcYAeAHEESwFPQQkaeFSBD4gLU0FFCwuQWJzdHJhY3R9fWENChkgTowEhQYEIC0BRmIfAAVaBashLg2qCFN1YgXQCGVzOhWoCR8OCQiR6wVfDaFBpBAJc3ViICVoAF8m6AkVQOGuAHg26AA5UipZCwQgLiG8DVxiIgCRQwR1Yjp6ACBUZXJtaW5hdGWWeQAuZQpaegA-IwA2ewBmZwA0CScgLS0tIE1lbWJlcnPRtkaJAQEmYS0Z7ARGSRJjDAGQARklDDpOAAgtLS0BG1o0AHwJJyBDYWxsZWQgYnkgQmFiZWwgcHJvdG9jb2wgdG8gdw4JCRggdGhpcyBvEhkMBA0KHvQJIcMWKAkUc18sIGpfwSbFJyhpXyA6IGlfID0gMGbzAC6GBxrCCQHFTSRaHgA2UAF6uwAMcmVhZGa6ABLWCA25lqYAIqAJWqgAMiAARqoAFG9udmVydC6RACFgBGEgDvMJ4UoIaW5nOVzVxQRUbwEgACglXAkMBCA9KacYVG9Kc29uKOUmBCwgDQcEbWUBKj4TByEpWn8AJG4gWE1MIGRvY3UO8wtSgQAIWE1MCYAOJQoEVG8BMBWDCFhtbBGC_QgELCBOiAABloUrYYc16ABpIvcJBCU-DRU="
	cCsharp_BaseCs = "zQa8e3tkZWZpbmUgIkNPTU1FTlRTIn19e3skY210cyA6PSBleHBhbmRDb21tZW50cyAuAR0YaWYgbGVuIAUkBRBwbmRlbnR9fS8vLyA8c3VtbWFyeT4NCnt7cmFuZ2UNNB0kEHt7Ln19ASAIZW5kQkAAAC8uQQANIQUHATERrBRTSU1QTEUysgBSbAA-agBCSQAUTUVUSE9EMkkABddO-wAxBUKzAC6yAAmHBUQ-KAA-jQAAaT0TADwy9wAJRSguUGFyYW1ldGVycz5LAAQ8cAEbSCBuYW1lPSJ7ey5OYW1lfX0iPnstcxAkaSwkbQHEAC4ytQAMZiAkaSl6LnwABHt7JRcFnw0MBDwvBWElxQ0dBQcuhAEMQVRUUk0tDGF0dHJFLhxmaWx0ZXJBdAEPMisCBSIuxgAAWx2rCCAkeAGsACQZJQmpBCwgDY8V4gR7e018MhwBACgVTBBqLCAkeQFMMiAABTYMJGp9fR1RCGlmIC02FV0EID0ZcCRmb3JtYXRWYWx1TcAF_wApFZAFDwBdIQ8FCg0Y"
	cCsharp_ClientCs = "9RPwPC8vIDxhdXRvLWdlbmVyYXRlZCAvPg0KLy8gQVVUTy1HRU5FUkFURUQgRklMRSAtIERPIE5PVCBNT0RJRlkFKABHFTnkZnJvbSB7ey5GaWxlbmFtZX19DQp7e3RlbXBsYXRlICJTSU1QTEVDT01NRU5UUyIgLkNvbW1lbnRzIAEqOA0KdXNpbmcgU3lzdGVtOzYPADAuQ29sbGVjdGlvbnMuBXIEaWM-IwAMTGlucT4UAAxUZXh0QhQANGhyZWFkaW5nLlRhc2tzFR8wQ29uY3VyLkJhYmVsOwG6FHJhbmdlIAWWCHN9fQUIAd8EfX0FHghlbmQBtwHpSHNwYWNlIHt7aW5kZXggLk5hbWUFEyRzICJjc2hhcnAiASoAexlUIC5TZXJ2aWNlcwUYDHtzZXQBOyBudCAiXHQifX0uNQFWLwEECVspJxwuQ29kZURvbSFLEHBpbGVyKR0hyAEaBCgiBc5UIiwgIiIpXQ0KCXB1YmxpYyBjbGFzcwHMAasofX1DbGllbnQgOiAFLgkOIEJhc2UsIEl7ewXPBH19LgwANEFzeW5jDQoJew0KCQkvQUMcc3VtbWFyeT4REQlEDWgoY29uc3RydWN0b3IRIAQ8Lz4yABg8cGFyYW0gIU8oPSJiYXNlVXJsIj4BhgQgcyktFCBVUkw8LwUnFWgyNgBAdGltZW91dFNlY29uZHMiPlQJEBAgaW4gcwkULj8ALQsV4QmzDChzdHJBhg2CFCwgaW50IDZZACQgPSAxMjApIDogASVEKG5ldyBIdHRwVHJhbnNwb3J0BRIlSkBKc29uU2VyaWFsaXplcigpLDYpAAAuQXIUYXRVcmwoAUoUVXJsLCAiFZAIIiksOnkAJCwgImFwcGxpY2FhCBAvanNvbiHaQhQAHCkpIHsgfQ0KNW0-bAEJ60aeASggdGhhdCBjYW4gYkH4SGVkIGZvciB1bml0IHRlc3RpbmcVU3bAAQB0MQwgIj5bTW9ja10gFRIYIG9iamVjdIqKAQBJJUQ1VxlDBCkgMXwVVgXyACAB82gjcmVnaW9uIFN5bmNocm9ub3VzIG1ldGhvZHN5gCwkaSwgJG0gOj0gLk0JHAB9hawyiQNCiwMUTUVUSE9EfZEIfX0JcVKUe3tmb3JtYXRUeXBlIC5SZXR1cm5zfX0ge3t0b1Bhc2NhbENhc2WJHBh9fSAoe3tyhV4FigB2BYoAUGEBEGV0ZXJzgQIuTwABVQVMbYwke3tpZiAuSW5pdE1bKH19ID0gbnVsbHt7hZcFIShsYXN0ICRpICRtLhlbFCB8IG5vdGHPFHt7ZWxzZQFrDTMFBwApIeJl2AAJBUUUaXNWb2lkHdIMU2VuZBE7AHIF50QgTWFrZVJlcXVlc3RBbmREZXNR8wA8WhYBAD4NpwAoXe4ELCBhSQREaaXkDGFyeTxpiQAsTTAEPigh4342ATJQAHVKBCB9BcuGEAGh0g2bDCB9KTshC0UXJRoUDQojZW5kSV4EDQpRZ6UL_mgCsmgCCGlmIEKVAcHJMZUBDH55AVqXAgXN_pwC_pwC6pwCTZBBowXLMRMNGGKoAgUtfjQB_q0C_q0C6q0CAAkuxgIIfQ0K"
	cCsharp_ModelCs = "mBfwPC8vIDxhdXRvLWdlbmVyYXRlZCAvPg0KLy8gQVVUTy1HRU5FUkFURUQgRklMRSAtIERPIE5PVCBNT0RJRlkFKABHFTnkZnJvbSB7ey5GaWxlbmFtZX19DQp7e3RlbXBsYXRlICJTSU1QTEVDT01NRU5UUyIgLkNvbW1lbnRzIAEqOHVzaW5nIFN5c3RlbTsNCi4PADAuQ29sbGVjdGlvbnMuBXAEaWM-IwAMVGV4dBUUMENvbmN1ci5CYWJlbDsBhRRyYW5nZSAFVAhzfX0FCAGqBH19BR4IZW5kAYIBtEhzcGFjZSB7e2luZGV4IC5OYW1lBRMkcyAiY3NoYXJwIgEqAHsZVBQuRW51bXMFFQx7c2V0ATggbnQgIlx0In19Lv0AVvcABAlbCfEcLkNvZGVEb20hExBwaWxlcgnnIZABGgQoIgXLTCIsICIiKV0NCglwdWJsaWMgZW51JXIBpwGDAAkdmTwkaSwkdiA6PSAuVmFsdWVzAY8caWYgJGl9fSwBvwX1DAkJe3sF5TB9fSA9IHt7Zm9ybWF0BTEEIC4BMgUmCA0KCSXLBQ0Ee3spTRQuQ29uc3T--gDe-gAsc2VhbGVkIGNsYXNzIcpOAgEV-QAJMTAAYwGzCCB7ewUIJFR5cGUgLkRhdGEBCgR9fRlHTgsBWSlKDAEcJGlzLCAkeHMhbxggLlN0cnVj5hoBLjYCDEFUVFJhKhxBdHRyaWJ1dCHHEc1Ae3tpZiAuQWJzdHJhY3R9fWENCggge3slqTocAQkvKEV4dGVuZHN9fSA6AfAVDwgsIElFYyhNb2RlbHt7ZWxzZQUkNhYAaShBZggJCS-BOyRzdW1tYXJ5Pg0KCREcRGVmYXVsdCAlYSEGBG9yERsEPC8uLQBNwVV-BCgpASgEeyBVXgxGaWVsAagJtyhJbml0aWFsaXplckXoQbgsdG9QYXNjYWxDYXNlaapNxRRjYXN0IC4p1TrTAjJIAAA7LSkJYwEwGC5Jc0xpc3R-YwAIbmV3dSwBNxFtBCgpUlMACE1hcP5SAAFSJdUhJgB9QiUBBA0KYdIyfwSagQSCbQJZbokcLvAACCB7e1KFARx7IGdldDsgcwEFAH0tZQGjdYgsb3ZlcnJpZGUgc3RyofYIVG9TBQlNBgHlKAl2YXIgc2VyID0gIW9FnRhKc29uU2VyTQkIKCk7BSmlzgAoBS8YdHJtID0gKK1UTElPLk1lbW9yeVN0cmVhbSlzZXIuFT8UKHRoaXMpRXZJ0CQJCXJldHVybiAoAXBcVVRGOEVuY29kaW5nKGZhbHNlKSkuR2V0Dag0c3RybS5Ub0FycmF5KCkNiQB9AY8FBRgjcmVnaW9uLlQDARxV9WG5dXs1DXF3GHZpcnR1YWxxyFR2b2lkIFJ1bk9uQ2hpbGRyZW48VD4oJRFlrgBB5RlAPFQ-IG1ldGhvZCwgVCBhdXiB2RwsIGJvb2wgcgE9JEFsbCA9IHRydWUF9QB7AZEMCWlmKAk1OCA9PSBudWxsKSB0aHJvdyV0DEFyZ3XhoCBOdWxsRXhjZXDhgQQoIgkxACIB9UKJAgAJYZlMaWYgaXNUcml2aWFsUHJvcGVydHlxIwhpZigRjgApEeFaRgIIPSAocb4udQIAKQ3nBCgilT0gIiwgdHlwZW9mVjEAACxaqAIALDEnCctlYIkMNZwlJRBiYXNlLkKBASkHFUQALDVmlQUBPQB9MvgCOl4AfvoBJb0ZeSH3bTEOqwgELCCSBAL-7gE57iQJCQlzd2l0Y2goAX8FWGE7TgUCBAljoZU9osFuAHRWsQVWtQH-5gFO5gEAIG3iYQX90AwJCQlkybgBpjKcAw0xOgICIYkhAwAsDaVVBSX75TsNOYU60RVBEAAJhSAFBQwjZW5kiSgyYAkUDQp9IA0K"
	cCsharp_MvcAsyncControllerCs = "sRHwPC8vIDxhdXRvLWdlbmVyYXRlZCAvPg0KLy8gQVVUTy1HRU5FUkFURUQgRklMRSAtIERPIE5PVCBNT0RJRlkFKABHFTmIZnJvbSB7ey5GaWxlbmFtZX19DQoNCnVzaW5nIFN5c3RlbTs2DwAwLkNvbGxlY3Rpb25zLgVIBGljFSNYQ29uY3VyLkJhYmVsOw0Ke3tyYW5nZSAFTwhzfX0FCAFuBH19BR4IZW5kAXABeEhzcGFjZSB7e2luZGV4IC5OYW1lBRMkcyAiY3NoYXJwIgEqBHsgEVNsJGssICRzIDo9IC5TZXJ2aWNlc319e3tpZiAkawUqAHsNXRB7e3NldAFWGG50ICJcdCIBJoR0ZW1wbGF0ZSAiQ09NTUVOVFMiIC5Db21tZW50cyB9fQlbCfscLkNvZGVEb20BHBBwaWxlcgnxIXIBGgQoIgXpdCIsICIiKV0NCglwdWJsaWMgcGFydGlhbCBjbGFzcwHvAc4YfX1Db250ciE9JHIgOiB7e2Jhc2UZExR9fTxJe3sF-Ch9fUFzeW5jPg0KCRnwHC5NZXRob2RzCdwwaWYgLlBhcmFtZXRlcgX5MuUABecECQk6gAAkUmVxdWVzdCA6IC6nARQuTXZjLkkFyQ0hCA0KCR15FCRpLCAkeCVpcnIABXQh84JdAQwJCXt7GSUMQVRUUiF_HEF0dHJpYnV0IcgtSDR7e2Zvcm1hdFR5cGUgLgEGQH19IHt7dG9QYXNjYWxDYXNlSSo2TwIoCQkJI3JlZ2lvbiBC2AAxo1R2b2lkIFJ1bk9uQ2hpbGRyZW48VD4oJQYUTW9kZWxBReEQPFQ-IG0leEgsIFQgYXV4RGF0YSwgYm9vbCByAT0oQWxsID0gdHJ1ZSkhNgQJeyEBSfQyoAEECQkB9UxpZiBpc1RyaXZpYWxQcm9wZXJ0eRHYCGlmKBFSCCkge0m9VusADCA9IChOGgEAKQ2rBCgiVUggIiwgdHlwZW9mVjEAACxaTQEALBHrACl5pggJCQmFGgQJCTFDJQY5QyFACHN0coE4YdAELCCSTQE1NyEiGHN3aXRjaCgBPwkYclABBAljQQAd7WFVThwC_jEBsjEBGCByZXR1cm5FFDI-ASQJCWRlZmF1bHQ6ESMUZmFsc2U7RTghXEZiAUWlDFNldEQJOQRzKCkWXmUChf4MLkluaYF9FGl6ZXJ9fSFWBGlmQTBWGQEUPSBudWxsQXFWIAAYIHt7Y2FzdFGlUXcUVmFsdWUgNmsAGe1NvgnQYY0IZW5kaZAB7AB9YSMAZan5OpwFhUMuHAQUTUVUSE9EvaQ2vQVOOwRVnyUOCGlzVmHpBC5SJZUIc319rdg4VGhyZWFkaW5nLlRhc2tzBQYUe3tlbHNlciMAADwxAYGMACAZVgA-DfO57CmrAHsB8QAJBYoyDwQQdmFyIHKpjoFSNCA9IERlc2VyaWFsaXplrYYAPHXGDREEPihhj4VJJT9NaQhtX2LhVyBlc3NMb2dpYy5WNQTFUgAo9T1e1AXhPwhpfX2BGQ1kxTMBpVpZAA36CZs53AQJfQ0ZEA0KfQ0K"
	cCsharp_MvcControllerCs = "7RDwPC8vIDxhdXRvLWdlbmVyYXRlZCAvPg0KLy8gQVVUTy1HRU5FUkFURUQgRklMRSAtIERPIE5PVCBNT0RJRlkFKABHFTmIZnJvbSB7ey5GaWxlbmFtZX19DQoNCnVzaW5nIFN5c3RlbTs2DwAwLkNvbGxlY3Rpb25zLgVIBGljFSNYQ29uY3VyLkJhYmVsOw0Ke3tyYW5nZSAFTwhzfX0FCAFuBH19BR4IZW5kAXABeEhzcGFjZSB7e2luZGV4IC5OYW1lBRMkcyAiY3NoYXJwIgEqBHsgEVNsJGssICRzIDo9IC5TZXJ2aWNlc319e3tpZiAkawUqAHsNXRB7e3NldAFWGG50ICJcdCIBJoR0ZW1wbGF0ZSAiQ09NTUVOVFMiIC5Db21tZW50cyB9fQlbCfscLkNvZGVEb20BHBBwaWxlcgnxIXIBGgQoIgXpdCIsICIiKV0NCglwdWJsaWMgcGFydGlhbCBjbGFzcwHvAc4YfX1Db250ciE9JHIgOiB7e2Jhc2UZExR9fTxJe3sF-BR9fT4NCgkZ6xwuTWV0aG9kcwnXMGlmIC5QYXJhbWV0ZXIF9DLgAAXiBAkJOnsAJFJlcXVlc3QgOiAuogEULk12Yy5JBcQNIQgNCgkdeRQkaSwgJHglZHJyAAV0Ie6CWAEMCQl7exklDEFUVFIhehxBdHRyaWJ1dCHDLUM0e3tmb3JtYXRUeXBlIC4BBkB9fSB7e3RvUGFzY2FsQ2FzZUklNkoCKAkJCSNyZWdpb24gQtgAMZ5Udm9pZCBSdW5PbkNoaWxkcmVuPFQ-KCUGFE1vZGVsQUXcEDxUPiBtJXhILCBUIGF1eERhdGEsIGJvb2wgcgE9KEFsbCA9IHRydWUpITYECXshAUnvMqABBAkJAfVMaWYgaXNUcml2aWFsUHJvcGVydHkR2AhpZigRUggpIHtJuFbrAAwgPSAoThoBACkNqwQoIlVDICIsIHR5cGVvZlYxAAAsWk0BACwR6wApeaEICQkJhRUECQkxQyUGOUMhQAhzdHKBM2HLBCwgkk0BNTchIhhzd2l0Y2goAT8JGHJQAQQJY0EAHe1hUE4cAv4xAbIxARggcmV0dXJuRRQyPgEkCQlkZWZhdWx0OhEjFGZhbHNlO0U4IVxGYgFFpQxTZXRECTkEcygpFl5lAoX5DC5JbmmBeBRpemVyfX0hVgRpZkEwVhkBFD0gbnVsbEFxViAAGCB7e2Nhc3RRpVF3FFZhbHVlIDZrABntTb4J0GGNCGVuZGmQAewAfWEjAGWp9DqXBYVDLhwEFE1FVEhPRL2fNrgFTjsEVZ8RrIE3CCAuUiWWAHOFQK1gKU4AewGUAAkpOy7gBBB2YXIgcqkxYfU0ID0gRGVzZXJpYWxpemWtKQA8dWkNEQQ-KGEyZewF4gVXCGlzVoFtEYMcIHwgbm90fX1NKC0tCG1fYuEYIGVzc0xvZ2ljLlb7AwAo1flelQXB-whpfX1h2gWCHcNaVAANgQm5OZ0ECX0NGRANCn0NCg=="
	cCsharp_ServiceCs = "hQrwPC8vIDxhdXRvLWdlbmVyYXRlZCAvPg0KLy8gQVVUTy1HRU5FUkFURUQgRklMRSAtIERPIE5PVCBNT0RJRlkFKABHFTnkZnJvbSB7ey5GaWxlbmFtZX19DQp7e3RlbXBsYXRlICJTSU1QTEVDT01NRU5UUyIgLkNvbW1lbnRzIAEqOA0KdXNpbmcgU3lzdGVtOzYPADAuQ29sbGVjdGlvbnMuBXIIaWM7AV4UcmFuZ2UgBToIc319BQgBgwR9fQUeCGVuZAFbAY1Ic3BhY2Uge3tpbmRleCAuTmFtZQUTJHMgImNzaGFycCIBKgx7IHt7CVNsJGssICRzIDo9IC5TZXJ2aWNlc319e3tpZiAkawUqAHsNXRB7e3NldAFWGG50ICJcdCIBJhn0Vu4ALhYBDEFUVFIhDRxBdHRyaWJ1dAFoBAlbKQYcLkNvZGVEb20hKhBwaWxlcgn8IacBGnQoIkJhYmVsIiwgIiIpXQ0KCXB1YmxpYyBpbnRlcmYB8AhJe3sF6wGyAAkZ3RwuTWV0aG9kcwEXOsAAQsIAFE1FVEhPRB3INuEAUr8AmAl7e2Zvcm1hdFR5cGUgLlJldHVybnN9fSB7e3RvUGFzY2FsQ2FzZSmHCH19KDV1EGksICR4JXUgUGFyYW1ldGVyNXcUaX19LCB7KXc2YAABZgVdDewNJQApWQQACUWL_qkB_qkBkqkBGEFzeW5jDQr-rgG2rgEgaWYgaXNWb2lkPa1NgzhUaHJlYWRpbmcuVGFza3MFBhx7e2Vsc2V9fWojAAA8NqMBVQMAPi2dWgsCBf7-EAKOEAINhxANCn0NCg=="
	cGo_BaseGo = "0gJ0e3tkZWZpbmUgIkNPTU1FTlRTIn19e3tyYW5nZSAuAQtQaW5kZW50fX0vL3t7Ln19Cnt7ZW5kARkFBwAKGUEUU0lNUExF_kcAAkcAFE1FVEhPRE5HABxDb21tZW50cwF9YpYADbogUGFyYW1ldGVyOi4AMCB7ey5OYW1lfX06IHsN5yQkaSwkbSA6PSAuMmIADGYgJGkF8BlsBCB7Ef0lGwBlIQk-EQE="
	cGo_InvokerGo = "yxPIe3t0ZW1wbGF0ZSAiU0lNUExFQ09NTUVOVFMiIC5Db21tZW50cyB9fQpwYWNrYWdlIHt7DQq0fX0KCi8vICoqKiBBVVRPLUdFTkVSQVRFRCBGSUxFIC0gRE8gTk9UIE1PRElGWQEoES9sR2VuZXJhdGVkIGZyb20ge3suRmlsZW5hbWV9fQUo7AppbXBvcnQgKHt7aWYgc2VydmljZVVzZXNUeXBlICJkZWNpbWFsIn19CgkibWF0aC9iaWcie3tlbmR9fVozABhhdGV0aW1lCTQFCg0wIAp7e3JhbmdlIAl0LHN9fQkie3sufX0iCg0iBCkKFSQsJGssICRzIDo9IC5TCZYAcw1yDCRrfX0RMRFSHC5NZXRob2RzBRpYc2V0aW5kZW50ICIifX0vLyB7eyRzLk4F9QR7ew0JPFJlcXVlc3QgaXMgdGhlIHINDyxzdHJ1Y3R1cmUgdXMhNihvciBpbnZva2luZwUoFUAEIG0Fcgggb24NGBVjMUYILgp0IURufQAJbhAge3sieyEvEc4UJGksICR4BfwgUGFyYW1ldGVyQtsACFx0IiERWVRWTgI8CXt7dG9QYXNjYWxDYXNlIA3_ICB7e2Zvcm1hdCXoAC4BBiB9fSBganNvbjohlTEuRHNlcmlhbGl6ZXJPcHRpb25zIA0sBCJgMYQACkmedEluaXQgc2V0cyBkZWZhdWx0IHZhbHVlcyBmb3IgYUF7CVQteSwKZnVuYyAob2JqICpqIwEEKSABWAQoKXIlADFBMUAALi42AUWmAC4BRg3JQcwMb2JqLloZATw9IG5ldyh7e25vdHB0ciAoPicBGCl9fSkKCSpyQgAxXQBWIQIAIDZ9AE3rCZcheBguSXNMaXN0gpcAEG1ha2UoEVoBOTGLDCwgMClOVgAITWFw_lUADVVl7ywKCXJldHVybiBvYmoh9AR7e5ppAxRzcG9uc2V5ag0QSf3-awOGawMybwBRKwAKJYAu3AEAUgXdVmcDAAkpyTZ9ARUxcTcFJSEWQjMDFSvqNgMJx4Y3AwkofjgDCSYxDmYNAQUaAC4NwDqqAgXwCCA9IEqZAg0xBH19QpwCDRoILklzXZ_CSgBakQIJF20MxkwGNtACLrAHVlwFAHQ64AUAIEnTMCB7CglTdmNPYmogSXvZe1EnDC0iYArBwekMaroGhoEAqUQAczIdBQApuWsMKHJlcTIcANW-rYUQLCByc3BeIgBJIxgpIGVycm9yAdRFDFIzAwRyZQktBCwgLXwIZXJywb8Ecy4pCAAuFXkAKLGkBCRw4eEEcHMBJ0awBQgkcGtBPA1OCHJlcT6uBQQkcPnSJfYAKWrSAwgJaWYBxhggPT0gbmlsAcsQCXJzcC5RdhG-CAoJfQ1xlesIZXJyJawFY9XxAAo="
	cGo_ModelGo = "8w_Ie3t0ZW1wbGF0ZSAiU0lNUExFQ09NTUVOVFMiIC5Db21tZW50cyB9fQpwYWNrYWdlIHt7DQq0fX0KCi8vICoqKiBBVVRPLUdFTkVSQVRFRCBGSUxFIC0gRE8gTk9UIE1PRElGWQEoES9sR2VuZXJhdGVkIGZyb20ge3suRmlsZW5hbWV9fQUo5AppbXBvcnQgKHt7aWYgbW9kZWxVc2VzVHlwZSAiZGVjaW1hbCJ9fQoJIm1hdGgvYmlnInt7ZW5kfX1SMQAYYXRldGltZQkyBQoNLiAKe3tyYW5nZSAJcCxzfX0JInt7Ln19IgoNIgQpChUkFC5FbnVtcwFjICRubSA6PSAuTgWyNHt7c2V0aW5kZW50ICIiASA5UlZMAQB0Acg0e3skbm19fSBzdHJpbmclSyhWYWx1ZXMgZm9yIC4jABwKY29uc3QgKBWKECRpLCR2BYIAVgUyCH19CQ1OBHt7DZUIID0gAcoJoi7OABQvLyBHZXQNLyggcmV0dXJucyB0aCHZCZFQZm9yIGEgZ2l2ZW4gaW50ZWdlciB2AWgULgpmdW5jHT8AKAUXASUAKSGrBT9ALCBib29sKSB7Cglzd2l0Y2gJPAQge26-AAhjYXMBgBhmb3JtYXRWAWwcIC59fToKCQkJoQAgDbAV3xQsIHRydWUxpRwJZGVmYXVsdB0xNCIiLCBmYWxzZQoJfQp9OV4oIGNvbnZlcnRzIGEtegggdG8F_wxlbnVtTZYhfQQncwHYLv0ADC8vIFIxLxxydWUgd2hlbgU6KcIEIGkhuhh1bmQgYW5kCXwJIwhub3QtOgwocyAqDcoEKSAF5xgoKSAoaW50PTUoaWYgcyA9PSBuaWwhRDEBADAdzzFZBCpzklYBJZc1Oj0kBHt7OnEBglYBLoYABH0KMX812hguQ29uc3RzYQP-IwM-8wIALiUvYaEtQRXRCCA9IELJABGaZb4AZYEYFaQcJGlzLCAkeHNhPxggLlN0cnVjAbIACuLGA3FxAHMFSxAge3sie4UFQGlmIC5FeHRlbmRzfX0KCXt7HQ6NXZF-EC5GaWVsBScyhgAEXHQFSX5OBCwJe3t0b1Bhc2NhbENhSBGRMRalNgAuAQYgfX0gYGpzb246mRsBcTxyaWFsaXplck9wdGlvbnMgDSwEImAxR21BJEluaXQgc2V0cyBtcGnbQfeBMjWSiScMKG9iakHvDY8EKSABPwQoKR0TACA1LKmFNQelxQAuATANmiFBEG9iai57VuoAPD0gbmV3KHt7bm90cHRyIChJ1wHyCcwYKX19KQoJKnJCADEuRWgAIDZ9ADWvIdQBVxguSXNMaXN0gpcAEG1ha2UoEVoBOQmQFH19LCAwKU5WAAhNYXD-VQANVUn2kT8Qb2JqCn0xzw=="
	cGo_ServiceGo = "wQXIe3t0ZW1wbGF0ZSAiU0lNUExFQ09NTUVOVFMiIC5Db21tZW50cyB9fQpwYWNrYWdlIHt7DQq0fX0KCi8vICoqKiBBVVRPLUdFTkVSQVRFRCBGSUxFIC0gRE8gTk9UIE1PRElGWQEoES9sR2VuZXJhdGVkIGZyb20ge3suRmlsZW5hbWV9fQUo7AppbXBvcnQgKHt7aWYgc2VydmljZVVzZXNUeXBlICJkZWNpbWFsIn19CgkibWF0aC9iaWcie3tlbmR9fVozABhhdGV0aW1lCTQFCg0wIAp7e3JhbmdlIAl0LHN9fQkie3sufX0iCg0iBCkKFSQsJGssICRzIDo9IC5TCZYAcw1yDCRrfX0RMTgKe3tzZXRpbmRlbnQgIiIBIjllWl8BJHR5cGUgSXt7Lk4pByxpbnRlcmZhY2UgeyARohwuTWV0aG9kcwVqHWIEXHQ6ZAAUTUVUSE9EHWoBgwmQAQowdG9QYXNjYWxDYXNlIA1xACgRZRQkaSwgJHgF4yBQYXJhbWV0ZXIV5RBpfX0sIC0XGawce3tmb3JtYXQlngAuAQYBZiWOBCkgJZAuIAAkUmV0dXJuc319KDY4ABUYPCwgZXJyb3Ipe3tlbHNlfX0FDgB7Kd4AIDFkAH0NEQAK"
	cJava_BaseJava = "ggTse3tkZWZpbmUgIkNPTU1FTlRTIn19Cnt7JGNtdHMgOj0gZXhwYW5kQ29tbWVudHMgLn19e3tpZiBsZW4gBSQFLkhpbmRlbnR9fS8qKgp7e3JhbmdlER4Ee3sRHRQgKiB7ey4FMAhlbmQFSA03ECAqL3t7DRQFBwAKGZcUU0lNUExFHZ0RYAmKDUIELy86WwA-RwAMQVRUUgHbGHt7JGF0dHIF2xxmaWx0ZXJBdAEPCU4N2AUiLqEAEXQUJGksICR4IRcJRgUkDGYgJGkF1hHpDcwcQHt7Lk5hbWUNJSEyOC5QYXJhbWV0ZXJzfX0oezEqEGosICR5AVYyIAAlaBQkan19LCANUgUSFVQNCQQgPRkfKGZvcm1hdFZhbHVlKaolBwApFT8FDwAKNg8A"
	cJava_ConstantsJava = "3gLwdS8vIEFVVE8tR0VORVJBVEVEIEZJTEUgLSBETyBOT1QgTU9ESUZZCi8vIEdlbmVyYXRlZCBmcm9tIHt7aWRsLkZpbGVuYW1lfX0Ke3tzZXRpbmRlbnQgIiJ9fXt7dGVtcGxhdGUgIlNJTVBMRUNPTU1FTlRTIiABO1RDb21tZW50cyB9fQpwYWNrYWdlIHt7DQoQfX07CgouQgAZPAAuMjkAPHVibGljIGNsYXNzIHt7Lk4FjAwgeyB7Lo4AWFx0In19CgkKe3tyYW5nZSAuVmFsdWVzAaQJsQh9fXAJSXRzdGF0aWMgZmluYWwge3tjb25zdFR5cGUgLkRhdGEBCgh9fSAZaCQ9IHt7Zm9ybWF0BVIsIC59fTsKe3tlbmR9AXAAfQ=="
	cJava_EnumJava = "-AfwdS8vIEFVVE8tR0VORVJBVEVEIEZJTEUgLSBETyBOT1QgTU9ESUZZCi8vIEdlbmVyYXRlZCBmcm9tIHt7aWRsLkZpbGVuYW1lfX0Ke3tzZXRpbmRlbnQgIiJ9fXt7dGVtcGxhdGUgIlNJTVBMRUNPTU1FTlRTIiABO1RDb21tZW50cyB9fQpwYWNrYWdlIHt7DQqIfX07CgppbXBvcnQgamF2YS5pby5TZXJpYWxpemFibGU7CgouYAAZWgAuMlcAIHVibGljIGVudQGwBC5OBakUIGltcGxlCXtcY29tLmNvbmN1ci5iYWJlbC5tb2RlbC5CAQwURW51bSwgLnUABCB7NuUABFx0BeccICRlIDo9IC4B8yxyYW5nZSAkaSwgJHYFFRhWYWx1ZXN9IR8pHAEmDYoMKHt7LgUgRH19KXt7aWYgbGFzdCAkaSAkZQkYXHMgfCBub3R9fSx7e2Vsc2V9fTt7e2VuZAFEBQchBhFXDHByaXYhbTRmaW5hbCBpbnQgdmFsdSkqFSQpEQEdDGdldFYBIEAoKSB7IHJldHVybiB0aGlzLgk2BCB9TlwABHt7EbsBQwUqACklFhFhBHt7EQoZTAQgPRGKGR42XgAJlhhzdGF0aWMgFWQYIGZpbmRCeQmmgnAAGHN3aXRjaCgRlEGgRYwhLBGUGYAEe3spqi1oAH02oAFGMAAIY2FzQZ4JLwx9fToKLksARi8AGR4tWhXVMcF6QQAYZGVmYXVsdMJoABBudWxsO1JcAC7iAAh9Cn0="
	cJava_InterfaceJava = "xibwdS8vIEFVVE8tR0VORVJBVEVEIEZJTEUgLSBETyBOT1QgTU9ESUZZCi8vIEdlbmVyYXRlZCBmcm9tIHt7aWRsLkZpbGVuYW1lfX0Ke3tzZXRpbmRlbnQgIiJ9fXt7dGVtcGxhdGUgIlNJTVBMRUNPTU1FTlRTIiABO1RDb21tZW50cyB9fQpwYWNrYWdlIHt7DQp0fX07CgppbXBvcnQgamF2YS51dGlsLkhhc2hNYXA7RhoABRYRF3Rjb20uY29uY3VyLmJhYmVsLlNlcnZpY2VNZXRob2QVPkInAABCAS0NLGYmABxSZXNwb25zZQ0pflUADFZvaWSaKwBQcHJvY2Vzc29yLkJhc2VJbnZva2VyZokAEHRyYW5zISIFLxRDbGllbnSOLgAAVBE4Zi0AAHMpMQAuLhIBWERlZmluaXRpb247Cnt7JHNydiA6PSAuJfoUcmFuZ2UgKbAIc319CQkcIHt7Ln19LioBLwhlbmQFKVkTWQ0ALjIKAih1YmxpYyBjbGFzcwFBAE5FXQFbBGxlSS9WlAAIIHsgMoECIFx0In19Cgp7e0mSEH19LyoqHQ5QICogR2V0cyB0aGUgaW50ZXJmYWNlDXYgZm9yIHRoaXMgMQYyOAAALx0OAHAJrQBDAa0EPEkBQxA-IGdldAUKBRUsKCkgeyByZXR1cm4gBRYALgXZBDsgdpsAMENvbnZlbmllbmNlIG1FtRggdG8gY3JlYU0UYSBuZXcgTR4UIGluc3RhASgkY29udGFpbmluZwXVDbwEIGkBljk_AGEhwAAuHcMEICpK0QAASUl7DXINDgAoLnkBJCBpRmFjZUltcGwd4h2YBCgoBe8AKRknfvQAAFQyigEAZEVtDcwAbSUDBHMgQpkBVCBZb3Ugc2hvdWxkIHByb3ZpZGUgYW4F7kEsBe0IIG9mKc415wQuIG2QQVwALgWoghIBFT0pxyAgZXh0ZW5kcyAyFAEAez1YWXkEe3tpHhQkaSwgJG1lNIk0AHM2JAVaEQMuSQBRzDh7e2Zvcm1hdFR5cGUgLlJFYEhzfX0ge3t0b0NhbWVsQ2FzZSAuAeAIfX0oMoEAabUgUGFyYW1ldGVyBYQuTQABU15KADx7e2lmIGxhc3QgJGkgJG0uGUY8IHwgbm90fX0sIHt7ZWxzZQHaZfMEe3sFBwApLgQEEcbBMxEMQp4Dic8YIGNhbiBiZU2JAGRhBhRtYWtlIGFxnBAgY2FsbEUZCGFueSnrcTMEcyBFRRBlZCBpbkkCbRSC7QEMc3RhdJWEDYk58QRzZQ0TCGltcJGURUUMIHsKCVICApUnoaIMKFN0cmGlBHVyZTsUc3VwZXIoAQ1hIU7tAV4-AAAsgbk8IHRpbWVvdXRJbk1pbGxpc2GOFVMALEIeAJZkANUhACDVYxlYGRNeTQAycAI68QJuIAH-1gL-1gLG1gJi8gEuyAEYdG9QYXNjYTKUA1G_iQkEID2l02G-RisAhr8DSvQDAHv2qgNxli7HAB0UDGYgaXMO1wide3H3zext9-FGAC5VMQQuacVXAChthwn8BCk7LqUHGXGBHYU7XTJ2NQfNpywgaXMgYSBjb21wb24Oiwpl_QRlIBIJCgQgZoH3HHdvcmsgdXNlhU-FQA2S7WwQcy4gIFQOAAgZRe1dFjAIES5OWwckIGFuZCBuZWVkc4GsRGJlIHJlZ2lzdGVyZWQgd2l0aOWczX9AUmVxdWVzdERpc3BhdGNoZXK2mQQR8C6aBA0UHowIAHt2kQQNLgAoib4N1_HUiUYuFQBh8D2OOayNogxNYXA8id0ELCAS_ggEPD_1KA3zKfQIPj4gDkoKMhQAAHMOKwlSXgAZaK5hAAhtYXBtVwBIFjUMnpsAACjFrF3uEtIMwdEZFGbRByBtYXAucHV0KCJhzUKLBwAi4SBOEAQWFwoEKTttOHoIAW1lAG0OCg05HhkKMigDGRctpymjDqgKLYEOJwghfh6pCgAiFt4LGj8IACIh9xlFGQoNTynqRUMW_QoeTgsuUwAFHwnlAU8d4AWsZlABHSoMcHJpdg4ACzKeByFnRjQFVXelCEKtBIG8DfxJkZG-UtsNADwOlQ4IcnNlDhoJnfIAPi2sBCB74fUdwkZiCR3FPT8Nzzb2BgFlcn8JBC5JDoQNKGFsaXplcn19ID0gEUAdGBhMaXRlcmFsDpwNDaVVWQGLFGYgbGVuICq-CTbdAlHnLeBWYgEh3QB9DV0BVEIlA3Y-ADIjCEYxARHUIQ4EIC4BBnIUAdaTCu29croBOa9OCgDF3yEgQukDJZgAdEaMCyIHC05PAGULMX4ZIS1K_rcDOrcDGU9pA0oFBErAAIUNbvABHE9iamVjdFtdgV8JWSYlDEE0UjwAGeyN2g5JDhVKDtkPcqUBZocBZf_SHgIEIH0hq0aNAQ6XCBEWAQwWowgAfQ=="
	cJava_ModelJava = "pxDwdy8vIEFVVE8tR0VORVJBVEVEIEZJTEUgLSBETyBOT1QgTU9ESUZZDQovLyBHZW5lcmF0ZWQgZnJvbSB7e2lkbC5GaWxlbmFtZX19DQp7e3NldGluZGVudCAiIn19e3t0ZW1wbGF0ZSAiU0lNUExFQ09NTUVOVFMiIAE8IENvbW1lbnRzIAE9JHBhY2thZ2Uge3sNCuB9fTsNCg0KaW1wb3J0IGNvbS5nb29nbGUuZ3Nvbi5hbm5vdGF0aW9ucy5TZXJpYWxpemVkTmFtZTsVNBhqYXZhLmlvFSAIYWJsAR4ce3tyYW5nZSAJWghzfX0JCSAge3sufX0uKjsByQhlbmQBlSR7eyRtZCA6PSAuCQ4Z0BnKAC4dxy7yAEBBVFRSUyIgLkF0dHJpYnV0ZQ3oWHVibGlje3tpZiAuQWJzdHJhY3R9fSBhDQsEe3sFexQgY2xhc3MBkwHTIUsBMChFeHRlbmRzfX0gZQkKASEVFRE7EGltcGxlKVkxGgH6CCB7CQHaCfwQLkZpZWwBSjK0AQhcdCKW6AABSSnpCH19QBFlKX8QKCJ7ey4JtQQiKS4oAAxwcml2QQs0e3tmb3JtYXRUeXBlIC4BBkB9fSB7e3RvQ2FtZWxDYXNlIA1AKSUMSW5pdCnfFHJ9fSA9IBFAHRgYTGl0ZXJhbCGmLQUJOwFdFC5Jc0xpcyFjED0gbmV3FT8BFQFDEYMEKClOOwAITWFwBXUuOgAITWFwUjkAADsNQUG2RlcBBHt7MTEAcEUCIbspJxAoKSB7fQk1GGlmIGxlbiBqngFqRwBR2BQkaSwgJHZFtjHmLogBGUIxSiEnMQdqigEgbGFzdCAkaSAkDZ8gIHwgbm90fX0sDfgNBzL1AQB7AXo6cwIBEzEMGY0QdGhpcy5SBwIl9EofAjlwGUYAfQ2FAWRKdwCdoEXsGTk1lVKdAhRnZXR0ZXIhr0FkIbEcIHJldHVybiBmtwAIOyB9hUuCdAAQdm9pZCBFMS5kAKadAQQpICVm_lMBJVMhOBFJAH0BDRBlbmR9fYU5gsgAIFN0cmluZyB0bwkJISYBQkbUAhFhOX4JMihCdWlsZGVyIHNiIGmsMhcAnVkgKCIpO3t7ICRzReWJsynuSvwCGWYZCiBzYi5hcHBlbmQBVUpbAgg6IiklvEaqABk6ZvoBDCArICKFgXUyAHM-MwMgIHt7ZWxzZX19kXYNcSVvGagZCk1vCHNiLjlgYQEFRSVn5XjFKg0KDTOBu2JJAQgpIikuUQABSQWDAQkR_Qx9DQp9"
	cJs_BaseJs = "mwbwPHt7ZGVmaW5lICJDT01NRU5UUyJ9fQ0Ke3skY210cyA6PSBleHBhbmRDb21tZW50cyAufX17e2lmIGxlbiAFJAkvKGluZGVudH19LyoqAT4QcmFuZ2UNLA0aFCAqIHt7LgkuCGVuZAVHFRsIL3t7DRQFBwQNCgFFEZoUU0lNUExFHaAEe3tKYAAELy8-XwBKSwAUTUVUSE9EMksABcVO6QAR8wmzUuEABTw2ygAJ5Q3PAcYNqyRQYXJhbWV0ZXJzOjAABEBwARoEIHsBOSROYW1lfX0gfSB7DeIQJGksJG0BnQAuGY4lfwQkaTaXAA1xJQMNDAF9JQMEe3sRvQggKi8VFjZkAQxBVFRSIfsYe3skYXR0ciX6HGZpbHRlckF0AQ8y9wEFIi7HAABbHbAIICR4AbEJRwUlBGYgAa4ELCANmAh7ey4J6AXJQUgALi4dAQAoFUwQaiwgJHkBTDIgAAU2DCRqfX0dUQhpZiAVVA0JBCA9GXAkZm9ybWF0VmFsdU2UJRcAKRWQBQ8AXTUaDRg="
	cJs_ClientJs = "1gfwPC8vIDxhdXRvLWdlbmVyYXRlZCAvPg0KLy8gQVVUTy1HRU5FUkFURUQgRklMRSAtIERPIE5PVCBNT0RJRlkFKABHFTnkZnJvbSB7ey5GaWxlbmFtZX19DQp7e3RlbXBsYXRlICJTSU1QTEVDT01NRU5UUyIgLkNvbW1lbnRzIAEqmA0KdmFyIEJBQkVMUlBDID0gcmVxdWlyZSgnYmFiZWxycGMnKTsNCgF0FRoMLi8nKQFkNHJhbmdlIHVzaW5nc319FR8Qe3sufX0FNpR7e2VuZH19e3skbnMgOj0gaW5kZXggLk5hbWVzcGFjZXMgImpzIgGHAYUMbnMgPRWKGC51dGlscy4B0QUrJChnbG9iYWwsICcFTgVhBHt7CYIgLlNlcnZpY2VzBWsEY2wFbAVmAVsuDAFeBgEBf8RjbGllbnQgPSBmdW5jdGlvbihiYXNlVXJsLCB0aW1lb3V0U2Vjb25kcywganNvblJwYyEFCA0KCTI9AAhuZXcZxABCIUMALglZhk4AADsNTAh0aGEBhxB0aGlzOyF6EecUJGksICRtIUgkLk1ldGhvZHN9fSGHCHNldCFZGG50ICJcdCIhBDn9FE1FVEhPRB33CH19CQFZOC57e3RvQ2FtZWxDYXNlIC0tLv0AMnYAAHYFdiBQYXJhbWV0ZXIlaUpAAAQsIC3yOGNhbGxiYWNrKXsNCgkJYyVaNC5zZW5kUmVxdWVzdCgiLcY4Lnt7JGNsc319IiwgInt7DYgYIiwgeyB7eyneBfdGgQAdLgA6Qf9BVQH8NGlmIGxhc3QgJGkgJG0uGbMcIHwgbm90fX0VpA2rDCB9LCARryAgKTsJCQkNCgklWABlQbwIDQp9IYQMbnNbJxWfECddID0gCdUEOyAhgg0q"
	cJs_ModelJs = "qArwPC8vIDxhdXRvLWdlbmVyYXRlZCAvPg0KLy8gQVVUTy1HRU5FUkFURUQgRklMRSAtIERPIE5PVCBNT0RJRlkFKABHFTnkZnJvbSB7ey5GaWxlbmFtZX19DQp7e3RlbXBsYXRlICJTSU1QTEVDT01NRU5UUyIgLkNvbW1lbnRzIAEqlA0KdmFyIEJBQkVMUlBDID0gcmVxdWlyZSgiYmFiZWxycGMiKQ0KAVI0cmFuZ2UgdXNpbmdzfX0RJyAne3sufX0nKTsBIwhlbmQZVAxucyA9FVkYLnV0aWxzLgGgbHNwYWNlKGdsb2JhbCwgJ3t7aW5kZXggLk5hbWUFHBRzICJqcyINVQFXCXocLkVudW1zfX0u3ABe1gAIbnNbAZkBUxx9fSddID0gexlJOCRpLCR2IDo9IC5WYWx1ZQVTHGlmICRpfX0sASYFxgwJInt7BZM0fX0iIDoge3tmb3JtYXQFMgQgLgGGBScEDQolbg0MBHt7DasUQ29uc3RzIQzSrgD-rQCyrQAcJGlzLCAkeHMhEhggLlN0cnVjprsALt8ADEFUVFJBYBxBdHRyaWJ1dCFaRokBMGZ1bmN0aW9uKCkNCnsuWgIULkZpZWxkJY5ZvVbhAX52AEB7e2lmIC5Jbml0aWFsaXplciGKFAl0aGlzLjXWQRQoe2Nhc3QgLlR5cGUh0jLjATI-ABw7e3tlbHNlIAFXAS4YLklzTGlzdFZXAARbXUovAAhNYXBWLgAEe30NLlYfAAxudWxsASFpWgR7e01qAAkJ4xBFeHRlbik0DGZ1bGxB8AxPZiAuFRcULmNhbGwoAfkAKRlFBA0KKQssdG9TdHJpbmcgPSBmNZVMew0KCQlyZXR1cm4gSlNPTi50b3MFJg1GDA0KCX1F9WEqIQ8VbJFEbesAZQWbACg6AQIALGFKTqoAHaAV5QQNCg=="
	cJs_PackageJson = "ngREewogICJuYW1lIjogIkxpYiIsAREYdmVyc2lvbgEUEDAuMC4wCRYcZGVzY3JpcHQNGlRHZW5lcmF0ZWQgZnJvbSB7ey5GaWxlAUwEfX0JMRRhdXRob3IBRgkQbGtleXdvcmRzIjogWwogICAgImJhYmVsIgogIF0FeDByZXBvc2l0b3J5IjogAZsYICAidHlwZQFDCGdpdAVGFCAgInVybAESHGh0dHBzOi8vAUAAfQVACGJ1ZwFdDTpiJwBAbWFpbiIgOiAibW9kZWwuanMFWQAiCd4dQRBzdGFydAFqJG5vZGUgaW5kZXgRLQGZBGVzBR0UbWFrZSB0AQ0ICiAgCYgoZGVwZW5kZW5jaWUdTxhuZXRtYXNrAVEUfjEuMC40BXccICAiYXN5bmMFFxQwLjIuMTARGCB1bmRlcnNjb3IlARBeMS42LhUcAGIhOgRycAU3AF4hpQAxOQcMZGV2REaCAAmbGGxpY2Vuc2UBbgVnFCJlbmdpbgUQBeYwPj0gdjAuMTEuKiIKfQ=="
	cJs_ServiceImplJs = "tAjwPC8vIDxhdXRvLWdlbmVyYXRlZCAvPg0KLy8gQVVUTy1HRU5FUkFURUQgRklMRSAtIERPIE5PVCBNT0RJRlkFKABHFTnkZnJvbSB7ey5GaWxlbmFtZX19DQp7e3RlbXBsYXRlICJTSU1QTEVDT01NRU5UUyIgLkNvbW1lbnRzIAEqkA0KdmFyIEJBQkVMUlBDID0gcmVxdWlyZSgnYmFiZWxycGMnKTsJJRRzZXJ2ZXIuIwAELi8BFBRpY2UuanMFJwF6NHJhbmdlIHVzaW5nc319FU-we3sufX0nKTt7e2VuZH19e3skeCA6PSBpbmRleCAuTmFtZXNwYWNlcyAianMiASAJSCAuU2VydmljZXMBrS7XAFbRABR7eyRjbHMBWQVTAUgIc2V0AWUYbnQgIlx0IgFHAfIJKUh9fSA9IGZ1bmN0aW9uKCkgew0KGcMUJGksICRtBUwUTWV0aG9kQoQAFE1FVEhPRB2KJH19CXRoaXMue3sNfy5fAAR7ewnTBVgAdgVYJFBhcmFtZXRlcnMBqih0b0NhbWVsQ2FzZSkSDH19LCAtMHxjYWxsYmFjayl7DQoJCXRocm93IG5ldyBFcnJvcignIhVyKCIgbm90IGltcGxpIeYEZWQlmQQJfTF5BA0KRSAlhQHuQG1vZHVsZS5leHBvcnRzID0gGQoBYDkdDHsgDQoJ3iBqc29ucnBjID0Fgi4iAErjAC7AAQQJIAlBFa0uHwEF2iQgCQlyZXR1cm4gQVoEZXI5SgAuDXYYLmNyZWF0ZQFXCGVyKAWCFU8QKCkgKTsBpwggCX0tOBANCgl9OxG5CGFwaf61ADq1ABhiYXNlVXJsfrwACGFwaTq4AA05ACySwQAltA=="
	cJs_ServiceJs = "qAzwPC8vIDxhdXRvLWdlbmVyYXRlZCAvPg0KLy8gQVVUTy1HRU5FUkFURUQgRklMRSAtIERPIE5PVCBNT0RJRlkFKABHFTnkZnJvbSB7ey5GaWxlbmFtZX19DQp7e3RlbXBsYXRlICJTSU1QTEVDT01NRU5UUyIgLkNvbW1lbnRzIAEqkA0KdmFyIEJBQkVMUlBDID0gcmVxdWlyZSgnYmFiZWxycGMnKTsJJRRqYXlzb24uIwAJEgUhAXQ0cmFuZ2UgdXNpbmdzfX0VSbB7ey59fScpO3t7ZW5kfX17eyR4IDo9IGluZGV4IC5OYW1lc3BhY2VzICJqcyIBIAlIIC5TZXJ2aWNlcwGnLtEAVssAFHt7JGNscwFZBVMBSAhzZXQBZRhudCAiXHQiAUcB7AkpTH19ID0gZnVuY3Rpb24oKSB7DQoJAR4sdGhhdCA9IHRoaXM7DRMYX2ltcGw7IBniFCRpLCAkbQVrFE1ldGhvZEKjABRNRVRIT0QdqQh9fQkBUwgue3sNni5-AAR7ewnyBVgAdgVYJFBhcmFtZXRlcnMBySh0b0NhbWVsQ2FzZSkxDH19LCAtTyBjYWxsYmFjaykBwgwJdHJ5BQgQCWlmICgFuQRbJxV0oCddID09PSB1bmRlZmluZWQpDQoJCQkJdGhyb3cgbmV3IEVycm9yKCciFTIUIiBub3QgIQAAaUEzBGVkJewICQkJBVsZzgQoIP7EAAXECCApOwGRAZccfWNhdGNoKGUJ1wAJEeZgKHsnY29kZSc6LTEsICdtZXNzYWdlJzplLg0LAH1BrQgJCX0BRxAJDQoJfVF2DA0KDQopgyBqc29ucnBjID0F81kDJTglpRRjcmVhdGVBfwRlci6oASEIDZBBKwBtKfwQID0ge31BLAQJCTEFAC5dFwgJCQkNKQRbJ0H7CH19LlGCOUAhmgAgQXsZEzK5AC1xBCA9JYwJbQGDDHMgPSBphRwuc2VydmVyKA1sKQccCXJldHVybiBFyyEVIRI5CghhcGnKBgEcYmFzZVVybCwFi_4PAWIPAWWEBGF0yhABkcwULkJhYmVsMRgVto4hAaVCha0hLTBtb2R1bGUuZXhwb3J0IXIZCmE1XUJOBgNBR5HGACBJgnVfLX4VEAgoKTtx1xQgDQp9DQo="
	cTest_BaseTest = "vAOoe3tkZWZpbmUgIklURU1TIn19IHt7JGluZCA6PSBnZXRpbmRlbnR9fQp7exELOCJpdGVtcyI6eyB7e2FkZAkXECAiICAiNicASHt7Z2V0VHlwZUtleSAuVmFsdWUBDix9fSI6Int7Zm9ybWEFHzYcAAx7e2lmHRJULklzQ29sbGVjdGlvbn19LHt7ZW5kfQWPEGYgYW5kXisALhgACEtleQl5HcQAawkUZosAAC4VNgAiNnUAGco6nAAge3t0ZW1wbGF0NU0dhQggfX0ZuwBzMVgAICFqARopRQB9BdkR4TGZMEJVSUxEX0NPTU1FTlQhogAKDUo="
	cTest_InterfaceTest = "wRmUeyB7e3NldGluZGVudCAiICAifX17eyRpZGwgOj0gaWRsfX0Ke3sJH_BafX0iX2NvbW1lbnQiOiJBVVRPLUdFTkVSQVRFRCBGSUxFIC0gRE8gTk9UIE1PRElGWS4gR2VuZXJhdGVkIGZyb20ge3tpZGwuRmlsZW5hbWV9fS4ge3tqb2luQwlWBHMgAXkALhEOFH19IiwgIC5_AGBzdHJ1Y3RzIjp7Cnt7cmFuZ2UgJGksICRzAa0MYWxsUwkgAcMRsx29DHt7Lk4FdwlACGlmIB1rGShOCgAAIhn4OrIAHUMYIix7e2VuZDIyAU5NAAR7ewF7MEV4dGVuZHN9fSJwYXIpUwh7ey4ZFppNACQicHJvcGVydGllRhcBAGYhFxAuRmllbAFoWosAJfouMgEAIuooARlk_jIBPjIBTqcALGdldFR5cGVLZXkgLgEJLH19Ijoie3tmb3JtYQUaFRcpczRJc0NvbGxlY3Rpb259fS5hARxmIC5Jc01hcC4IARnMdgoADCJrZXkBfgV8BYgELksJElFvQmoALocAOoADACANAWWIPHRlbXBsYXRlICJJVEVNUyIJ0AAgAaZGtQOaewItSQB9JRMkbGFzdCAkaSAkc00nFCB8IG5vdCkdCGxzZQFqEVgRCEZaAzkgLi0EHRYybAB5qKptAAR9LC4IBAxlbnVtRu8CAGVB7wxhbGxFAR4BukaqAO7cAgAi_gQELgQEECJ2YWx1SrMDAHYBxAQuVgUdVsMAOVcZCjLXAAB7CT8FPgRmIDXiAGUJFwBzqnUBGWEZCpriASWEXmsARicFNt8BwR1pTq2wJf1arwXilwEoIm1ldGhvZHMiOls2TgYAbSWEAE0JHgB9VvEAwYGhxtmsOSNOCgDipgBSTQBCIQFhHTxmIC5IYXNQYXJhbWV0ZXJzRQFt7k5GAAgsInABLEZIAwBwJQCmSQAZXTKIAgmUBCRw_SIZKp4KAGY4ATJZAKp7B0ISAQgie3v-SQb-SQb-SQb-SQb-SQbdSQRtLlk5XgAEAApGAQROEAJhetV1AGbBmDwgLlJldHVybnMuSXNWb2lkhuMCAHIJMBryCUL8AU53ADoQAg1ySlwIGRrF4B2cAENmZwgNNb5vCDZlCBFFdmgIESY29whScAhabggNQupxCAB9busBOoQIAC6tc143AkZLAAhdCn0="
)

var staticFiles = map[string]string{
	"/error.babel": c_ErrorBabel,
	"/asp/_base.asp": cAsp_BaseAsp,
	"/asp/client.asp": cAsp_ClientAsp,
	"/asp/model.asp": cAsp_ModelAsp,
	"/csharp/_base.cs": cCsharp_BaseCs,
	"/csharp/client.cs": cCsharp_ClientCs,
	"/csharp/model.cs": cCsharp_ModelCs,
	"/csharp/mvcAsyncController.cs": cCsharp_MvcAsyncControllerCs,
	"/csharp/mvcController.cs": cCsharp_MvcControllerCs,
	"/csharp/service.cs": cCsharp_ServiceCs,
	"/go/_base.go": cGo_BaseGo,
	"/go/invoker.go": cGo_InvokerGo,
	"/go/model.go": cGo_ModelGo,
	"/go/service.go": cGo_ServiceGo,
	"/java/_base.java": cJava_BaseJava,
	"/java/constants.java": cJava_ConstantsJava,
	"/java/enum.java": cJava_EnumJava,
	"/java/interface.java": cJava_InterfaceJava,
	"/java/model.java": cJava_ModelJava,
	"/js/_base.js": cJs_BaseJs,
	"/js/client.js": cJs_ClientJs,
	"/js/model.js": cJs_ModelJs,
	"/js/package.json": cJs_PackageJson,
	"/js/service-impl.js": cJs_ServiceImplJs,
	"/js/service.js": cJs_ServiceJs,
	"/test/_base.test": cTest_BaseTest,
	"/test/interface.test": cTest_InterfaceTest,
}

func Lookup(path string) []byte {
	s, ok := staticFiles[path]
	if !ok {
		return nil
	} else {
		d, err := base64.URLEncoding.DecodeString(s)
		if err != nil {
			log.Print("babeltemplates.Lookup: ", err)
			return nil
		}
		r, err := snappy.Decode(nil, d)
		if err != nil {
			log.Print("babeltemplates.Lookup: ", err)
			return nil
		}
		return r
	}
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	if strings.HasSuffix(p, "/") {
		p += "index.html"
	}
	b := Lookup(p)
	if b != nil {
		mt := mime.TypeByExtension(path.Ext(p))
		if mt != "" {
			w.Header().Set("Content-Type", mt)
		}
		w.Write(b)
	} else {
		http.NotFound(w, r)
	}
}
