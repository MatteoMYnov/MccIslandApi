window.addEventListener("scroll", function () {
    const scrollPosition = window.scrollY;
    const maxScroll = document.body.scrollHeight - window.innerHeight;

    let opacity = 0;
    if (scrollPosition > (maxScroll * 0.3)) {
        opacity = Math.min((scrollPosition - (maxScroll * 0.30)) / (maxScroll * 0.30), 1);
    }

    document.body.style.background = `rgba(23, 23, 23, ${opacity})`;

    const playerCard = document.getElementById("playerCard");
    const hypixelStatsCard = document.getElementById("hypixelStatsCard");

    if (scrollPosition > 0) {
        playerCard.style.top = 20 + scrollPosition * -0.05 + "%";
        minecraftStatsCard.style.top = 40 + scrollPosition * -0.05 + "%";
        hypixelStatsCard.style.top = 80 + scrollPosition * -0.05 + "%";
    }
});