//===================================================================
//  WEBCHAT.JS
//  provides a clean way to communicate with the websocket server
//===================================================================

var nickname = "nobody"
Webchat = {

    connect: function (host) {
        ws           = new WebSocket(host);
        ws.onopen    = Webchat.wsopen;
        ws.onclose   = Webchat.wsclose;
        ws.onerror   = Webchat.wserror;
        ws.onmessage = Webchat.wsmessage;
    },


    //===================================================================
    //  WEBSOCKET
    //===================================================================

    wsopen: function () {
        console.log("connected");
        Webchat.send("nickname",nickname)
    },

    wsclose: function () {
        console.log("disconnected")
    },

    wserror: function (e) {
        console.log(e)
    },

    wsmessage: function (e) {
        var x = JSON.parse(e.data);
        switch( x.Action )
        {
            case "inform":
                Webchat.oninform(x.Data);
                break;

            case "message":
                Webchat.onmessage(x.Origin, x.Data);
                break;
        }
    },


    //===================================================================
    //  MISC.
    //===================================================================

    compatible: function() {
        return window.WebSocket && window.localStorage;
    },

    //===================================================================
    //  API
    //===================================================================

    send: function (action, data) {
        ws.send(JSON.stringify( { Action: action, Data: data}));
    },

    oninform:   function () {console.log("`oninform` not implemented")},
    onmessage:  function () {console.log("`onmessage` not implemented")},

    //===================================================================
    //  LOCALSTORAGE
    //===================================================================

    load: function () {
        nickname = localStorage["nickname"]
        $('#settings-nickname').val(nickname)
    },
    save: function () {
        nickname = $('#settings-nickname').val()
        Webchat.send("nickname", nickname)
        localStorage["nickname"] = nickname;
    }
    
} // Webchat.JS
