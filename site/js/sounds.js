const uiclickSounds = [
    "/sounds/ui-click_normal-1.mp3",
    "/sounds/ui-click_normal-2.mp3",
    "/sounds/ui-click_normal-3.mp3"
];

function playUiClickSound() {
    const randomIndex = Math.floor(Math.random() * uiclickSounds.length);
    const sound = new Audio(uiclickSounds[randomIndex]);
    sound.volume = 0.2; // Volume réduit à 10%
    sound.play();
}

window.playUiClickSound = playUiClickSound;
