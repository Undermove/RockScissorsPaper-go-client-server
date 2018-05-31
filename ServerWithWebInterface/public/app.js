new Vue({
    el: '#app',

    data: {
        ws: null, // Our websocket
        newMsg: '', // Holds new messages to be sent to the server
        chatContent: '', // A running list of chat messages displayed on the screen
        email: null, // Email address used for grabbing an avatar
        username: null, // Our username
        joined: false, // True if email and username have been filled in
        rooms: '',
        newRoom: ''
    },

    created: function() {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.addEventListener('message', function(e) {
            var msg = JSON.parse(e.data);
            if (msg.type == 'message') {
                self.chatContent += '<div class="chip">'
                + '<img src="' + self.gravatarURL(msg.email) + '">' // Avatar
                + msg.username
                + '</div>'
                + emojione.toImage(msg.message) + '<br/>'; // Parse emojis
                var element = document.getElementById('chat-messages');
                element.scrollTop = element.scrollHeight;
            } else if(msg.type == 'createRoom'){
                var msg = JSON.parse(e.data);
                self.rooms += '<li class="collection-item"><div>'+msg.message+'<a href="#!" class="secondary-content"><i class="material-icons">meeting_room</i></a></div></li>'
    
                var element = document.getElementById('rooms-list');
                element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
            }
        });
    },
    
    methods: {
        send: function () {
            if (this.newMsg != '') {
                var a = {
                    type: "message",
                    email: this.email,
                    username: this.username,
                    message: $('<p>').html(this.newMsg).text() // Strip out html
                }
                this.ws.send(
                    JSON.stringify(a));
                this.newMsg = ''; // Reset newMsg
            }
        },

        createRoom: function () {
            if (this.newRoom != '') {
                this.ws.send(
                    JSON.stringify({
                        type: 'createRoom',
                        email: this.email,
                        username: this.username,
                        message:$('<p>').html(this.newRoom).text()
                    }
                ));
                this.newRoom = ''; // Reset newMsg
            }
        },

        join: function () {
            if (!this.email) {
                Materialize.toast('You must enter an email', 2000);
                return
            }
            if (!this.username) {
                Materialize.toast('You must choose a username', 2000);
                return
            }
            if (!this.newRoom) {
                Materialize.toast('You must enter room name', 2000);
                return
            }
            this.email = $('<p>').html(this.email).text();
            this.username = $('<p>').html(this.username).text();
            this.joined = true;
            this.newRoom = $('<p>').html(this.newRoom).text()
        },

        gravatarURL: function(email) {
            return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
        }
    }
});