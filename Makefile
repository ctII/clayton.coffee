all:
	hugo --minify
	rm -f public/assets/*.js
	rm -f public/assets/pink.css
	rm -f public/assets/red.css
	rm -f public/assets/blue.css
