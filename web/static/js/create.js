const inoph = document.getElementById("inoph");
const iulph = document.getElementById("iulph");
const ise = document.getElementById("ise");

const payload = document.getElementById("payload");

function updateSyntax() {
    const text = payload.value;
    const re = /{{(.*?)}}/gm;

    const placeholders = [];

    // update <li>
    iulph.innerHTML = "";
    let m;
    do {
        m = re.exec(text);
        if (m) {
            const ph = m[1].trim();

            if (!placeholders.includes(ph)) {
                placeholders.push(ph);

                iulph.innerHTML += `
<li class="list-group-item d-flex justify-content-between lh-condensed bg-dark">
  <div>
    <h6 class="my-0 bg-dark text-light">${ph}</h6>
    <small class="text-muted">${m[0]}</small>
  </div>
</li>
`;
                // build json data
                if (ph.startsWith(".")) {
                    const spl = ph.split(".");
                    for (let i = spl.length - 1; i >= 1; i--) {
                    }
                }
            }
        }
    } while (m);

    let open = 0;
    let closed = 0;

    let lastOpen = false;
    let lastClosed = false;

    for (let s of text) {
        if (s === "{") {
            lastClosed = false;
            if (lastOpen) {
                lastOpen = false;
                open++;
            } else {
                lastOpen = true;
            }
        } else if (s === "}") {
            lastOpen = false;
            if (lastClosed) {
                lastClosed = false;
                closed++;
            } else {
                lastClosed = true;
            }
        } else {
            lastOpen = false;
            lastClosed = false;
        }
    }

    if (open !== closed) {
        ise.innerHTML =
            "Probably a syntax error: " +
            open +
            " open <-> " +
            closed +
            " closed.";
        ise.style.display = "flex";
    } else {
        ise.style.display = "none";
        inoph.innerHTML = placeholders.length.toString();
    }
}

updateSyntax();
payload.onkeyup = updateSyntax;