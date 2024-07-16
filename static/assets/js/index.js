let launcherBtn = document.querySelectorAll("#launcher");
let launcherLogsClear = document.getElementById("launcherLogsClear");
let editBtn = document.querySelectorAll("#edit");
let winetricksBtn = document.getElementById("winetricks");
let winetricksVerbs = document.getElementById("winetricksVerbs");
let winetricksLogsClear = document.getElementById("winetricksLogsClear");
let saveWinetricksVerb = document.getElementById("winetricksVerbsSave");

async function runFetch(gameId, ele) {
  if (gameId) {
    let eventSource = new EventSource(`/run/${gameId}`);

    eventSource.onmessage = (event) => {
      let logger = document.getElementById("launcherLogging");
      logger.scrollTop = logger.scrollHeight;

      if (event.data != "0") {
        logger.value += `${event.data}\n`;
      }

      if (event.data == "0") {
        eventSource.close();
        ele.textContent = "Launch";
        ele.removeAttribute("disabled");
      }
    };
  }
}

async function runWinetricks(gameId, verbs) {
  if (gameId && verbs) {
    let eventSource = new EventSource(`/winetricks/${gameId}/${verbs}`);

    eventSource.onmessage = (event) => {
      let logger = document.getElementById("commandLogs");
      logger.scrollTop = logger.scrollHeight;

      if (event.data != "0") {
        logger.value += `${event.data}\n`;
      }

      if (event.data == "0") {
        eventSource.close();
        winetricksBtn.innerText = "Winetricks";
        saveWinetricksVerb.innerText = "Save";
        saveWinetricksVerb.removeAttribute("disabled");
      }
    };
  }
}

launcherBtn.forEach((ele, _) => {
  ele.addEventListener("click", async () => {
    ele.textContent = "Process is running...";
    ele.setAttribute("disabled", true);

    let gameId = ele.dataset.gameId;

    await runFetch(gameId, ele);
  });
});

launcherLogsClear.addEventListener("click", () => {
  document.getElementById("launcherLogging").value = "";
});

editBtn.forEach((ele, _) => {
  ele.addEventListener("click", () => {
    gameId = ele.dataset.gameId;

    window.location.href = `/edit/${gameId}`;
  });
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

winetricksLogsClear.addEventListener("click", () => {
  document.getElementById("commandLogs").value = "";
});
