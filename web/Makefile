.PHONY: tailwindcss
tailwindcss:
	@npx @tailwindcss/cli -i internal/presentation/css/styles.css -o static/css/styles.css --watch

.PHONY: templ
templ:
	@templ generate view -watch
	