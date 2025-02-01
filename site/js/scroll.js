window.addEventListener("scroll", function () { 
    const scrollPosition = window.scrollY;
    const maxScroll = document.body.scrollHeight - window.innerHeight;

    let opacity = 0;
    if (scrollPosition > (maxScroll * 0.1)) {
        opacity = Math.min((scrollPosition - (maxScroll * 0.10)) / (maxScroll * 0.10), 1);
    }
    const bgdark = document.getElementById("bgdark");
    if (bgdark) {
        bgdark.style.opacity = opacity;
    }
});