function saveTreeState() {
  document.querySelectorAll("details[data-key]").forEach(d => {
    localStorage.setItem("tree:" + d.dataset.key, d.open ? "1" : "0");
  });
}

function restoreTreeState() {
  document.querySelectorAll("details[data-key]").forEach(d => {
    const state = localStorage.getItem("tree:" + d.dataset.key);
    if (state === "1") d.open = true;
  });
}

/* save when toggling */
document.body.addEventListener("toggle", e => {
  if (e.target.matches("details[data-key]")) saveTreeState();
});

/* restore after sidebar swap */
document.body.addEventListener("htmx:afterSwap", e => {
  if (e.target.id === "sidebar") restoreTreeState();
});

/* initial load restore */
window.addEventListener("load", restoreTreeState);