const form = $("#cmsAttribForm");
const data = $("#cmsPageDataForm");

const api = "/cms/create/req";

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

    const payload = {
        ...data.serializeFormJSON(),
        ...form.serializeFormJSON()
    }
    console.log("payload:", payload);

    Swal.queue([
        {
            position: 'top-end',
            title: "Create CMS Page " + payload.page_url,
            confirmButtonText: "Create",
            showLoaderOnConfirm: true,
            preConfirm: async () => {

                // make request
                let json = JSON.stringify(payload);
                let res = await fetch(api, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: json
                });

                // get response
                let text = await res.text();
                if (res.status === 401 && text === "") {
                    text = "Unauthorized, you are not allowed to create a cms page";
                }

                if (res.status !== 200) {
                    return Swal.insertQueueStep({
                        position: 'top-end',
                        icon: "error",
                        title: "Error:",
                        text: text
                    });
                }

                const html = `
                        <h2>Page created successfully!</h2>
                        
                        <h3>Payload:</h3>
                        <pre class="language-json">${json}</pre>
                        
                        <h3>Response</h3>
                        <code>${text}</code>`;

                return Swal.fire({
                    position: 'top-end',
                    icon: "success",
                    title: "Result:",
                    html: html,

                    showCancelButton: true,
                    cancelButtonText: "🌱 Add another page",
                    confirmButtonText: "👉 Visit page"
                }).then((result) => {
                    if (result.isConfirmed) {
                        window.location = payload.page_url;
                    } else {
                        window.location = `/cms/create`;
                    }
                });
            }
        }
    ])
})