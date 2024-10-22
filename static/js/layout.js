// The main JavaScript file that's paired with _layout.html

function configureUnpoly() {
	// Automatically use unpoly throughout the site
	up.link.config.followSelectors.push("a[href]");
	up.link.config.preloadSelectors.push("a[href]");
	up.form.config.submitSelectors.push(["form"]);

	// Allow forms to handle other method types
	// this is the first part of the Triptych Proposals
	// https://alexanderpetros.com/triptych/
	up.on("up:form:submit", (event) => {
		const method = event.form.getAttribute("method") || "GET";

		if (method == "GET" || method == "POST") {
			return;
		}

		event.preventDefault();

		const url = event.form.getAttribute("action");
		const opts = { method };

		if (method !== "DELETE") {
			opts.body = event.params.toFormData();
		} else if (formData != null) {
			url = event.params.toUrl(url);
		}

		fetch(url, opts).then(console.log);
	});
}

// This code is added in a `defer` so we are assuming the eager script tags
// before it would have already loaded if they could
const __cdnLoaded = typeof up !== "undefined";

// Behold, my cunning plan
// If we fail to load from CDN then inject a local version
// All the benefits of CDNs but with a fallback on local connections
// or for the more overzealous adblocker

if (__cdnLoaded) {
	configureUnpoly();
} else {
	console.warn("Unpoly did not load from CDN, using local fallback");

	// This uuh... this takes a minute to work.
	// Depending on your browser the timeout for the initial CDN failure can be
	// really long. I've seen Firefox take a few *minutes* to fail out.
	// But it does eventually run this code, always.

	let script = document.createElement("script");
	script.src = "/static/vendor/unpoly.min.js";
	script.defer = true;
	document.head.appendChild(script);

	let style = document.createElement("link");
	style.rel = "stylesheet";
	style.href = "/static/vendor/unpoly.min.css";
	document.head.appendChild(style);

	console.log("vendored dependencies injected, waiting...");

	// Yeah I know, but we're already in this deep
	setTimeout(() => {
		console.log("configuring unpoly and calling up.boot()");
		configureUnpoly();
		up.boot();
	}, 200);
}
