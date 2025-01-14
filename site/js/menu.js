function showPlayer() {
    var playerName = document.getElementById("playerInput").value;
    var playerCard = document.getElementById("playerCard");
    var playerNameElement = document.getElementById("playerName");


    var hypixelStatsCard = document.getElementById("hypixelStatsCard");
    var hypixelStatsElement = document.getElementById("hypixelStats");

    if (playerName) {
        playerCard.style.display = "block";
        hypixelStatsCard.style.display = "block";

        playerNameElement.textContent = playerName;

        hypixelStatsElement.textContent = "Voici les statistiques Hypixel du joueur " + playerName;

    } else {
        alert("Please enter a player name.");
    }
}

document.getElementById("playerInput").addEventListener("keydown", function(event) {
    if (event.key === "Enter") {
        showPlayer();
    }
});

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

window.onload = function() {
    var defaultPlayerName = "Leroidesafk";

    var playerCard = document.getElementById("playerCard");

    var hypixelStatsCard = document.getElementById("hypixelStatsCard");
    var hypixelStatsElement = document.getElementById("hypixelStats");

    playerCard.style.display = "block";
    minecraftStatsCard.style.display = "block";
    hypixelStatsCard.style.display = "block";

    playerNameElement.textContent = defaultPlayerName;

    hypixelStatsElement.textContent = "Voici les statistiques Hypixel du joueur " + defaultPlayerName;
};
