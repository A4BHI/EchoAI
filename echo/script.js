async function send() {
  let msg = document.getElementById("inp").value;
  let msgbox = document.createElement("div");
  msgbox.className = "message";

  let user = document.createElement("div");
  user.className = "user";
  user.innerText = msg;

  Object.assign(user.style, {
    backgroundColor: "#3b82f6",
    color: "#ffffff",
  });

  msgbox.appendChild(user);
  document.querySelector(".chat-content").appendChild(msgbox);

  let aitypebox = document.createElement("div");
  aitypebox.className = "message";
  let aitype = document.createElement("div");
  aitype.className = "ai";
  setTimeout(() => {
    Object.assign(aitype.style, {
      backgroundColor: "#09ba6a",
      color: "#ffffff",
      border: "2px solid rgb(106, 60, 60)",
    });
    aitype.innerText = "Thinking....";
    aitypebox.appendChild(aitype);
    document.querySelector(".chat-content").appendChild(aitypebox);
  }, 500);

  let state = true;
  let intervalid = setInterval(() => {
    aitype.innerText = state ? "Thinking...." : "Thinking....|";
    state = !state;
  }, 500);

  try {
    let response = await fetch("/send", {
      method: "POST",
      headers: {
        "content-type": "text/plain",
      },
      body: msg,
    });

    let data = await response.text();
    clearInterval(intervalid);
    if (data != null) {
      aitype.innerText = data;
    } else {
      aitype.innerText = "No response";
    }
  } catch (err) {
    clearInterval(intervalid);
    aitype.innerHTML =
      "ERROR: Internal Error Occurred. <br> Try After Sometime";
    console.error(err);
  }
}

