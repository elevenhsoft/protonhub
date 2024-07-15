let launcherBtn = document.getElementById("launcher");
let editBtn = document.getElementById("edit");

async function runFetch(gameId) {
  const response = await fetch("/run", {
    method: "POST",
    body: JSON.stringify({ gameId: gameId }),
  });

  return response.status;
}

launcherBtn.addEventListener("click", async () => {
  launcherBtn.textContent = "Process is running...";
  launcherBtn.setAttribute("disabled", true);

  let gameId = launcherBtn.dataset.gameId;

  let status = await runFetch(gameId);

  if (status === 200) {
    launcherBtn.textContent = "Launch";
    launcherBtn.removeAttribute("disabled");
  }
});

editBtn.addEventListener("click", () => {
  gameId = editBtn.dataset.gameId;

  window.location.href = `/edit/${gameId}`;
});
