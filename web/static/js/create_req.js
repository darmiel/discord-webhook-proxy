const form = $("#createForm");
const data = $("#createDataForm");

const api = "/dashboard/create/req";

(function ($) {
    $.fn.serializeFormJSON = function () {
        const o = {};
        const a = this.serializeArray();
        $.each(a, function () {
            if (o[this.name]) {
                if (!o[this.name].push) {
                    o[this.name] = [o[this.name]];
                }
                o[this.name].push(this.value || '');
            } else {
                o[this.name] = this.value || '';
            }
        });
        return o;
    };
})(jQuery);

// save last data
let lastData = `{
  "Camera": {
    "ID": "camera-1",
    "Name": "Camera 1",
    "Avatar": "https://source.unsplash.com/random/400x400"
  }
}`;

form.on("submit", (event) => {
    event.preventDefault();

    Swal.queue([
        {
            position: 'top-end',
            input: 'textarea',
            inputPlaceholder: "{\n\n}",
            inputValue: lastData,
            title: "Send Test Webhook",
            confirmButtonText: "Test & Save Webhook",
            text: "We will now send your payload to Discord to test if it is accepted. " +
                "If this works, your webhook will be saved. " +
                "If this does not work, please check your input again.\n" +
                "If your payload makes use of templates outside of a value, you can still specify sample data here that will be used for checking.",
            showLoaderOnConfirm: true,
            preConfirm: async (args) => {
                lastData = args;

                // create payload
                const p = data.serializeFormJSON();
                p["args"] = args; // append example args

                const json = JSON.stringify(p);

                const f = await fetch(api, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: json
                });

                // get text
                const text = await f.text();

                if (f.status !== 200) {
                    return Swal.insertQueueStep({
                        position: 'top-end',
                        icon: "error",
                        title: "Error:",
                        text: text
                    });
                } else {
                    // decode response
                    const resp = JSON.parse(text);

                    const {webhook, sent_json: sentJson} = resp;
                    const {uid, user_id: userId, secret} = webhook;

                    const url = `/call/json/${userId}/${uid}/${secret}`;

                    const html = `
                        <h2>Webhook created successfully!</h2>
                        <ul style="list-style-type: none; padding: 0; margin: 0;">
                            <li><strong>UID:</strong></li>
                            <li>${uid}</li>
                            <li><strong>Secret:</strong></li>
                            <li>${secret}</li>
                        </ul>
                        
                        <h3>Call URL</h3>
                        <input id="swal-input1" class="swal2-input" readonly value="${url}">
                        
                        <h3>Request Data</h3>
                        <pre class="lang-json" style="text-align: left !important;">${sentJson}</pre>`;

                    return Swal.fire({
                        position: 'top-end',
                        icon: "success",
                        title: "Result:",
                        html: html,

                        showCancelButton: true,
                        cancelButtonText: "🌱 Add another webhook",
                        confirmButtonText: "👉 View/Edit webhook"
                    }).then((result) => {
                        if (result.isConfirmed) {
                            window.location = `/dashboard/${uid}`;
                        } else {
                            window.location = `/dashboard/create`;
                        }
                    })
                }
            }
        }
    ])
})