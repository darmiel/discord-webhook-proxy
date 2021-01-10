const form = $("#createForm");
const data = $("#createDataForm");

const api = "/create/req";

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

form.on("submit", (event) => {
    event.preventDefault();

    Swal.queue([
        {
            title: "Send Test Webhook",
            confirmButtonText: "Test & Save Webhook",
            text: "We will now send your payload to Discord to test if it is accepted. " +
                "If this works, your webhook will be saved. " +
                "If this does not work, please check your input again.",
            showLoaderOnConfirm: true,
            preConfirm: async () => {
                // create payload
                const json = JSON.stringify(data.serializeFormJSON());
                console.log("Json Payload:", json);

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
                        icon: "error",
                        title: "Error:",
                        text: text
                    });
                } else {
                    return Swal.insertQueueStep({
                        icon: "success",
                        title: "Result:",
                        text: text
                    })
                }
            }
        }
    ]);
})