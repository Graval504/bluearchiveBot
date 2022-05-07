# Bluearchive Bot
## What is this?
This is Packge for making my discord bot.\
It helps users to access information of bluearchive from namu.wiki with commands.\
It is my first project, and yet work in process, so only few functions are provided.\
## To use this package
You can install this package using\
`go get github.com/Graval504/bluearchiveBot`

If it doesn't work, you can\
`go env -w GO111MODULE=auto`

## Requirement for using this package
you need python and beautifulsoup4 library because it uses python to send request.\
because sending request with go causes 403 error so can't use request from net/http package.
## How To Use (Example)
* Getting Bluearchive data with jsonfile
    ```go
    package main

    import (
	"github.com/Graval504/bluearchiveBot"
    )

    func main() {
    	data := bluearchiveBot.GetCharacterInfoFromData(bluearchiveBot.GetCharacterList())
    	bluearchiveBot.CreateJsonFileFromData(data)
    }
```

It creates jsonfile.json on your directory.\
this file provides data of bluearchive that this package will using.
* If you get error, you should visit namu.wiki website and solve captcha.
