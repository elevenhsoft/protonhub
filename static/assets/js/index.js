let launcherBtn = document.getElementById("launcher");
let editBtn = document.getElementById("edit");
let winetricksBtn = document.getElementById("winetricks");
let winetricksVerbs = document.getElementById("winetricksVerbs");
let saveWinetricksVerb = document.getElementById("winetricksVerbsSave");

async function runFetch(gameId) {
  const response = await fetch("/run", {
    method: "POST",
    body: JSON.stringify({ gameId: gameId }),
  });

  return response.status;
}

async function runWinetricks(gameId, verbs) {
  if (gameId && verbs) {
    let eventSource = new EventSource(`/winetricks/${gameId}/${verbs}`);

    eventSource.onmessage = (event) => {
      let logger = document.getElementById("commandLogs");

      if (event.data != "0") {
        logger.value += `${event.data}\n`;
      }

      if (event.data == "0") {
        eventSource.close()
        winetricksBtn.innerText = "Winetricks";
        saveWinetricksVerb.innerText = "Save";
        saveWinetricksVerb.removeAttribute("disabled");
      }
    };

  }
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

saveWinetricksVerb.addEventListener("click", async () => {
  let gameId = winetricksBtn.dataset.gameId;
  let verbs = winetricksVerbs.value;

  if (gameId && verbs) {
    winetricksBtn.innerText = "Process is running...";
    saveWinetricksVerb.innerText = "Process is running...";
    saveWinetricksVerb.setAttribute("disabled", true);
  }

  await runWinetricks(gameId, verbs);
});
