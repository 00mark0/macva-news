// Font Loading API to prevent layout shifts
if ('fonts' in document) {
    // Load font faces before rendering content
    const fontFaces = [
        new FontFace('Merriweather', 'url(https://fonts.gstatic.com/s/merriweather/v31/u-4n0qyriQwlOrhSvowK_l52_wFZVcf6.woff2)', {
            weight: '400'
        }),
        new FontFace('Merriweather', 'url(https://fonts.gstatic.com/s/merriweather/v31/u-4m0qyriQwlOrhSvowK_l5-eRZOf-I.woff2)', {
            weight: '700'
        })
        // Add any other font weights/styles you're using
    ];

    // Load all fonts in parallel
    Promise.all(fontFaces.map(font => font.load()))
        .then(fonts => {
            fonts.forEach(font => document.fonts.add(font));
            // Add a class once fonts are loaded to trigger rendering
            document.documentElement.classList.add('fonts-loaded');
        })
        .catch(error => {
            console.error('Font loading failed:', error);
            // Allow page to render anyway
            document.documentElement.classList.add('fonts-loaded');
        });
}
