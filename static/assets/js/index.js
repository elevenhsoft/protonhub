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
  try {
    const response = await fetch("/winetricks", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ gameId: gameId, verbs: verbs }),
    });
    if (!response.ok) {
      throw new Error("Network response was not ok " + response.statusText);
    }

    return await response.json();
  } catch (e) {
    console.error(e);
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

saveWinetricksVerb.addEventListener("click", async (e) => {
  e.preventDefault();

  winetricksBtn.innerText = "Process is running...";
  saveWinetricksVerb.innerText = "Process is running...";
  saveWinetricksVerb.setAttribute("disabled", true);

  let gameId = winetricksBtn.dataset.gameId;
  let verbs = winetricksVerbs.value;

  let logs = await runWinetricks(gameId, verbs);

  let logger = document.getElementById("commandLogs");
  logger.value = logs["log"];
  winetricksBtn.innerText = "Winetricks";
  saveWinetricksVerb.innerText = "Save";
  saveWinetricksVerb.removeAttribute("disabled");
});
