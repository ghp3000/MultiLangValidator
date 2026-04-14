

``` go
	v := validator.NewMultiLangValidator(validator.LocaleZh)
	if err := v.Register(validator.LocaleZh, "zh.txt"); err != nil {
		fmt.Println(err)
	}
	if err := v.Register(validator.LocaleEn, ""); err != nil {
		fmt.Println(err)
	}
	if err := v.Register(validator.LocaleZhTw, "zh.txt"); err != nil {
		fmt.Println(err)
	}
	type RegistrationForm struct {
		Username string `validate:"required,min=3,max=20" label:"用户名"`
		Email    string `validate:"required,email" label:"邮箱"`
		Password string `validate:"required,min=8" label:"密码"`
	}
	form := RegistrationForm{
		Username: "ab",
		Email:    "invalid",
		Password: "123",
	}
	fmt.Println(v.Validate(&form, "zh"))
	fmt.Println(v.Validate(&form, "en"))
	fmt.Println(v.Validate(&form, "zh_tw"))
```
zh.txt文件内容
```
RegistrationForm.Username=用户名
RegistrationForm.Email=邮箱
RegistrationForm.Password=密码
```
输出结果:
```
用户名长度必须至少为3个字符
Username must be at least 3 characters in length
用户名長度必須至少為3個字元
```