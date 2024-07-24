package browser

const tmpl = `<html>
	<head>
		<title>reboxed</title>
		<style>
			body {margin: 0px; font-family: Helvetica; background-color: #36393D; color: #EEE;}
			a {color: #FFF; text-decoration: none;}
			a:hover {color: #0AF;}
			.nav {padding: 8px; background-color: #4096EE; height: 20px; border-bottom: 1px solid #90C6FE; box-shadow: 0px 16px 16px rgba(0, 0, 0, 0.1);}
			.nav a {margin: 10px; font-size: 20px; font-weight: bolder;}
			.logo h1 {margin: 0px; font-size: 20px; font-style: italic; float: right; color: #FFF;}
			.items {padding: 16px 8px; height:80%; overflow: auto;}
			.item {margin-left: 2px; margin-right: 2px; display: inline-block; font-size: 11px; font-weight: bolder; width: 128px; height: 125px; text-align: center; text-shadow: 1px 1px 1px #000; text-overflow: ellipsis; overflow: hidden; white-space: nowrap; letter-spacing: -0.1px;}
			.item img {width: 128px; height: 100px;}
			.thumb {background-position: center;}
			.pagenav {padding: 16px 8px; float: right;}
			.pagenav a {margin: 8px; font-weight: bolder;}
		</style>
	</head>
	<body>
		<div class="nav">
			<div class="logo"><h1>reboxed</h1></div>
			{{if .InGame}}
				{{if not (.Category | eq "map")}}
					<a href="/browse/entities">Entities</a>
					<a href="/browse/weapons">Weapons</a>
					<a href="/browse/props">Props</a>
					<a href="/browse/saves">Saves</a>
				{{end}}
			{{else}}
				<a href="/browse/entities">Entities</a>
				<a href="/browse/weapons">Weapons</a>
				<a href="/browse/props">Props</a>
				<a href="/browse/saves">Saves</a>
				<a href="/browse/maps">Maps</a>
			{{end}}
		</div>
		<div class="items">
			{{range .Packages}}
				<div class="item">
					<a href="garrysmod://{{if (.Type | eq "map")}}install{{else}}spawn{{end}}/{{.Type}}/{{.ID}}/{{.Revision}}">
						<div class="thumb" style="background-image: url(//image.reboxed.fun/{{.ID}}_thumb_128.png), url(//image.reboxed.fun/no_thumb_128.png);">
							<img src="//image.reboxed.fun/overlay_128.png">
						</div>
						{{.Name}}
					</a>
				</div>
			{{end}}
		</div>
		<div class="pagenav">
			<a href="{{.PrevLink}}">Previous</a>{{.PageNum}}<a href="{{.NextLink}}">Next</a>
		</div>
	</body>
</html>`
