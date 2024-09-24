package main

import (
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/buke/quickjs-go"
	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/lmittmann/tint"
)

const port = ":4000"

type application struct {
	BackendBundle string
	ClientBundle  string
	rt            *quickjs.Runtime
}

var textEncoderPolyfill = `function TextEncoder(){} TextEncoder.prototype.encode=function(string){var octets=[],length=string.length,i=0;while(i<length){var codePoint=string.codePointAt(i),c=0,bits=0;codePoint<=0x7F?(c=0,bits=0x00):codePoint<=0x7FF?(c=6,bits=0xC0):codePoint<=0xFFFF?(c=12,bits=0xE0):codePoint<=0x1FFFFF&&(c=18,bits=0xF0),octets.push(bits|(codePoint>>c)),c-=6;while(c>=0){octets.push(0x80|((codePoint>>c)&0x3F)),c-=6}i+=codePoint>=0x10000?2:1}return octets};function TextDecoder(){} TextDecoder.prototype.decode=function(octets){var string="",i=0;while(i<octets.length){var octet=octets[i],bytesNeeded=0,codePoint=0;octet<=0x7F?(bytesNeeded=0,codePoint=octet&0xFF):octet<=0xDF?(bytesNeeded=1,codePoint=octet&0x1F):octet<=0xEF?(bytesNeeded=2,codePoint=octet&0x0F):octet<=0xF4&&(bytesNeeded=3,codePoint=octet&0x07),octets.length-i-bytesNeeded>0?function(){for(var k=0;k<bytesNeeded;){octet=octets[i+k+1],codePoint=(codePoint<<6)|(octet&0x3F),k+=1}}():codePoint=0xFFFD,bytesNeeded=octets.length-i,string+=String.fromCodePoint(codePoint),i+=bytesNeeded+1}return string};`
var processPolyfill = `var process = {env: {NODE_ENV: "production"}};`
var consolePolyfill = `var console = {log: function(){}};`

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>React App</title>
</head>
<body>
    <div id="root">{{.RenderedContent}}</div>
	<script type="module">
	  {{ .JS }}
	</script>
	<script>window.APP_PROPS = {{.InitialProps}};</script>
</body>
</html>
`

type PageData struct {
	RenderedContent template.HTML
	InitialProps    template.JS
	JS              template.JS
}

type InitialProps struct {
	Name          string
	InitialNumber int
}

func buildBackend() string {
	result := esbuild.Build(esbuild.BuildOptions{
		EntryPoints:       []string{"../frontend/serverEntry.tsx"},
		Bundle:            true,
		Tsconfig:          "../frontend/tsconfig.json",
		Write:             false,
		LogLevel:          esbuild.LogLevelError,
		Outdir:            "/",
		Format:            esbuild.FormatIIFE,
		Platform:          esbuild.PlatformNode,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Target:            esbuild.ES2015,
		Banner: map[string]string{
			"js": textEncoderPolyfill + processPolyfill + consolePolyfill,
		},
		Loader: map[string]esbuild.Loader{
			".tsx": esbuild.LoaderTSX,
		},
	})

	script := fmt.Sprintf("%s", result.OutputFiles[0].Contents)
	return script
}

func buildClient() string {
	clientResult := esbuild.Build(esbuild.BuildOptions{
		EntryPoints:       []string{"../frontend/clientEntry.tsx"},
		Bundle:            true,
		Write:             true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		LogLevel:          esbuild.LogLevelError,
	})
	clientBundleString := string(clientResult.OutputFiles[0].Contents)
	return clientBundleString
}

var APP_ENV string

func main() {

	app := &application{}
	app.BackendBundle = buildBackend()
	app.ClientBundle = buildClient()

	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
	slog.SetDefault(logger)
	rt := quickjs.NewRuntime()
	defer rt.Close()

	app.rt = &rt
	srv := &http.Server{
		Addr:              port,
		Handler:           app.routes(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      30 * time.Second,
	}

	slog.Debug("START", slog.String("port", port))
	slog.Info("START", slog.String("message", fmt.Sprintf("Server is running on port %s", port)))
	slog.Info("START", slog.String("url", fmt.Sprintf("http://localhost%s", port)))

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
