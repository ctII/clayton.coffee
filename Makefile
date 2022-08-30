all:
	hugo --minify --baseURL "${{ steps.pages.outputs.base_url }}/"
	rm -f public/assets/*.js
	rm -f public/assets/pink.css
	rm -f public/assets/red.css
	rm -f public/assets/blue.css
