var project_name="WebRepo";

//Functions for toast messages

async function sleep(time) {
    return new Promise(resolve => setTimeout(resolve, time));
}

function ToastIcon(icon, text) {
    M.toast({unsafeHTML: "<span class=\"material-icons\">"+icon+"</span>"+ text});
}

function ToastError(text) {
    ToastIcon("error", "Error (see console): "+text);
}

//Functions for using backend


(function(e,t) {
    window.addEventListener('load', init);
    // https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Proxy
    // Pretty fast way to get url params
    const params = new Proxy(new URLSearchParams(window.location.search), {
        get: (searchParams, prop) => searchParams.get(prop),
      });
    function init() {
        M.AutoInit();
    }
})();
