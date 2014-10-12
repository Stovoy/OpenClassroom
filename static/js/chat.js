oc = oc || {};

oc.chat = oc.chat || {};

var identifier = Math.random();
var $dataStore;
var id;
var $chatInput;
var $chatButton;
var $chatHistory;
var lastMessageID = 0;
var colors = {};
var rcolor = new RColor;

$(document).ready(function() {
    // Load datastore
    $dataStore = $('#data-store');
    id = $dataStore.data('page');

    setInterval(oc.chat.refresh, 1000);
    oc.chat.refresh();

    $chatInput = $('#chat-input');
    $chatButton = $('#chat-button');

    $chatHistory = $('#chat-history');

    $chatInput.on('keydown', function(e) {
        if (e.keyCode == 13) { // Enter
            oc.chat.sendMessage();
        }
    });

    $chatButton.on('click', function() {
        oc.chat.sendMessage();
    })
});

oc.chat.sendMessage = function() {
    var message = $chatInput.val();
    $.ajax({
        method: "POST",
        url: "/chat/message/",
        data: {
            page: id,
            message: message}
    }).done(function(data) {
        data = JSON.parse(data);
        if (data.Error) {
            console.log(data);
        } else {
            $chatInput.val("");
            oc.chat.refresh();
        }
    });
};

oc.chat.refresh = function() {
    $.ajax({
        url: "/chat/loadNew/",
        data: {
            page: id,
            lastMessage: lastMessageID,
            identifier: identifier}
    }).done(function(data) {
        data = JSON.parse(data);
        if (data.Error) {
            console.log(data);
        }
        var i;
        if (data.Users) {
            var users = data.Users;
            users.sort(function(a,b){
                if (a.Name == b.Name) return 0;
                if (a.Name < b.Name) return -1;
                if (a.Name > b.Name) return 1;
            });
            var $usersHeader = $('#users-header');
            var $chatUserArea = $('#chat-user-area');
            $usersHeader.siblings().remove();
            for (i = 0; i < users.length; i++) {
                var color = "l2";
                if (i % 2 == 0) {
                    color = "l1";
                }
                var type = "guest-line";
                var name = users[i].Name;
                if (!users[i].IsGuest) {
                    type = "user-line";
                    name =
                        '<a href="/user/' + name +
                            '">' + name + '</a>'
                }
                $chatUserArea.append(
                        '<span class="' +
                        type + ' ' +
                        color + '">' +
                        name +
                        '</span>');
            }
        }
        if (data.NewMessages) {
            var messages = data.NewMessages;
            for (i = 0; i < messages.length; i++) {
                if (parseInt(messages[i].ID) > parseInt(lastMessageID)) {
                    lastMessageID = messages[i].ID;
                }
                var username = messages[i].User;
                var color = colors[username];
                if (color === undefined) {
                    color = rcolor.get(true);
                    colors[username] = color;
                }
                var timestamp = '[' + messages[i].Time + ']';
                var user =
                    '<a href="/user/' + username + '">' +
                    username + '</a>';
                var $message = $('<div class="message">' +
                    timestamp + ' ' + user + ': ' +
                    messages[i].Message +
                    '</div>');
                $message.css('background-color', color);
                $chatHistory.append($message);
            }
        }
    });
};