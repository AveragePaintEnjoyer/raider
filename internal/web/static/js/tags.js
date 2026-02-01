function initTagSelect(container) {
  const input = container.querySelector(".tag-search");
  const dropdown = container.querySelector(".tag-dropdown");
  const selected = container.querySelector("#selected-tags");
  const hidden = container.querySelector("#tags-input");

  function updateHidden() {
    const ids = [...selected.querySelectorAll(".tag-chip")]
      .map(c => c.dataset.id);
    hidden.value = ids.join(",");
  }

  // Open dropdown
  container.querySelector(".tag-input")
    .addEventListener("click", () => {
      dropdown.style.display = "block";
      input.focus();
    });

  // Close on outside click
  document.addEventListener("click", e => {
    if (!container.contains(e.target)) {
      dropdown.style.display = "none";
    }
  });

  // Add tag
  dropdown.addEventListener("click", e => {
    const opt = e.target.closest(".tag-option");
    if (!opt) return;

    const id = opt.dataset.id;
    const name = opt.textContent.trim();

    if (selected.querySelector(`[data-id="${id}"]`)) return;

    const chip = document.createElement("span");
    chip.className = "tag-chip";
    chip.dataset.id = id;
    chip.innerHTML = `${name} <button type="button" class="remove">&times;</button>`;
    selected.appendChild(chip);

    updateHidden();
  });

  // Remove tag
  selected.addEventListener("click", e => {
    if (e.target.classList.contains("remove")) {
      e.target.closest(".tag-chip").remove();
      updateHidden();
    }
  });

  updateHidden();
}

document.body.addEventListener("htmx:afterSwap", e => {
  e.target.querySelectorAll?.(".tag-select").forEach(initTagSelect);
});