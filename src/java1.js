
let genreSelect = document.getElementById("genre2");
if (localStorage.getItem("lastGenre")) {
    genreSelect.value = localStorage.getItem("lastGenre");
    }
    genreSelect.addEventListener("change", function() {
    localStorage.setItem("lastGenre", genreSelect.value);
});