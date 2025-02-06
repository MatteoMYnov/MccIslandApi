window.addEventListener("scroll", function () { 
    const scrollPosition = window.scrollY;
    const maxScroll = document.body.scrollHeight - window.innerHeight;

    let opacity = Math.min(scrollPosition / (maxScroll * 0.5), 1);
    let blurValue = Math.min(scrollPosition / (maxScroll * 0.3) * 5, 5); // Flou max de 10px

    const bgdark = document.getElementById("bgdark");
    if (bgdark) {
        bgdark.style.opacity = opacity;
    }

    // Appliquer le flou uniquement sur le fond via la variable CSS
    document.documentElement.style.setProperty("--blur-value", `${blurValue}px`);
});
