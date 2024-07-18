(async () => {
  let launcherBtn = document.querySelectorAll("#launcher");
  let stopLauncherBtn = document.querySelectorAll(".stopLauncher");

  await fetch("/running-games")
    .then((response) => {
      return response.json();
    })
    .then((data) => {
      if (data.ids != null) {
        for (let i = 0; i < data.ids.length; i++) {
          let gameId = data.ids[i];

          const targetElement = Array.from(launcherBtn).find(
            (element) => element.dataset.gameId === gameId,
          );

          targetElement.textContent = "Process is running...";
          targetElement.setAttribute("disabled", true);

          const stopTarget = Array.from(stopLauncherBtn).find(
            (element) => element.dataset.gameId === gameId,
          );

          stopTarget.removeAttribute("disabled");
        }
      }
    });
})();

(async () => {
  let launcherBtn = document.querySelectorAll("#launcher");
  let stopLauncherBtn = document.querySelectorAll(".stopLauncher");
  let deleteLauncherBtn = document.querySelectorAll(".deleteLauncher");
  let launcherLogsClear = document.getElementById("launcherLogsClear");
  let editBtn = document.querySelectorAll("#edit");
  let winetricksBtn = document.getElementById("winetricks");
  let winetricksVerbs = document.getElementById("winetricksVerbs");
  let winetricksLogsClear = document.getElementById("winetricksLogsClear");
  let saveWinetricksVerb = document.getElementById("winetricksVerbsSave");

  async function runFetch(gameId, ele) {
    if (gameId) {
      let eventSource = new EventSource(`/run/${gameId}`);

      eventSource.onmessage = async (event) => {
        let logger = document.getElementById("launcherLogging");
        logger.scrollTop = logger.scrollHeight;

        if (event.data.includes("pid:")) {
          pid = event.data.replace(/^pid: /, "");

          const stopTarget = Array.from(stopLauncherBtn).find(
            (element) => element.dataset.gameId === gameId,
          );

          stopTarget.removeAttribute("disabled");
        }

        logger.value += `${event.data}\n`;

        if (event.data.includes("exit: 0")) {
          ele.textContent = "Launch";
          ele.removeAttribute("disabled");

          eventSource.close();
        }
      };
    }
  }

  async function stopFetch(gameId) {
    if (gameId) {
      await fetch(`/stop/${gameId}`);

      const stopTarget = Array.from(stopLauncherBtn).find(
        (element) => element.dataset.gameId === gameId,
      );

      stopTarget.setAttribute("disabled", true);
    }
  }

  async function deleteFetch(gameId) {
    if (gameId) {
      await fetch(`/delete`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ "gameid": gameId })
      });
    }

    location.reload();
  }

  async function runWinetricks(gameId, verbs) {
    if (gameId && verbs) {
      let eventSource = new EventSource(`/winetricks/${gameId}/${verbs}`);

      eventSource.onmessage = (event) => {
        let logger = document.getElementById("commandLogs");
        logger.scrollTop = logger.scrollHeight;

        logger.value += `${event.data}\n`;

        if (event.data.includes("exit: 0")) {
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
      let confirmation = confirm("Are you sure?");

      if (confirmation == true) {
        ele.textContent = "Process is running...";
        ele.setAttribute("disabled", true);

        let gameId = ele.dataset.gameId;

        await runFetch(gameId, ele);
      }
    });
  });

  stopLauncherBtn.forEach((ele, _) => {
    ele.addEventListener("click", async () => {
      let gameId = ele.dataset.gameId;

      await stopFetch(gameId);
    });
  });


  deleteLauncherBtn.forEach((ele, _) => {
    ele.addEventListener("click", async () => {
      let confirmation = confirm("Are you sure?");

      if (confirmation == true) {
        let gameId = ele.dataset.gameId;

        await deleteFetch(gameId);
      }
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
})();
