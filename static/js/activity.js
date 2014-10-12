oc = oc || {};

oc.activity = oc.activity || {};
var $globalActivity;
var $userActivity;

$(document).ready(function() {
    $globalActivity = $('#global-activity');
    $userActivity = $('#user-activity');
    if ($globalActivity.length > 0) {
        oc.activity.loadGlobalActivity($globalActivity);
    }
    if ($userActivity.length > 0) {
        var username = $('#data-store').data('user');
        oc.activity.loadUserActivity(username, $userActivity);
    }
});

oc.activity.loadGlobalActivity = function($globalActivity) {
    $.ajax({
        url: "/activity/global/"
    }).done(function(data) {
        data = JSON.parse(data);
        if (data.Error) {
            console.log(data);
        }
        var i;
        if (data.Activities) {
            var activities = data.Activities;
            $globalActivity.empty();
            for (i = 0; i < activities.length; i++) {
                var a = activities[i];
                var user = '<a class="user-link" href="/user/' + a.User + '">' +
                    a.User + '</a>';
                var timestamp = '[' + a.Time + ']';
                var action = a.Action;
                var re = /(\/.*\/)(.*)/gi;
                action = action.replace(re, '<a href="$1$2">$2</a>');
                $globalActivity.append('<div class="activity-item">' +
                    timestamp + " " +
                    user + ": " +
                    action +
                    '</div>');
            }
        }
    });
};

oc.activity.loadUserActivity = function(username, $userActivity) {
    $.ajax({
        url: "/activity/" + username + "/"
    }).done(function(data) {
        data = JSON.parse(data);
        if (data.Error) {
            console.log(data);
        }

        var i;
        if (data.Activities) {
            var activities = data.Activities;
            $userActivity.empty();
            for (i = 0; i < activities.length; i++) {
                var a = activities[i];
                var user = '<a class="user-link" href="/user/' + a.User + '">' +
                    a.User + '</a>';
                var timestamp = '[' + a.Time + ']';
                var action = a.Action;
                var re = /(\/.*\/)(.*)/gi;
                action = action.replace(re, '<a href="$1$2">$2</a>');
                $userActivity.append('<div class="activity-item">' +
                    timestamp + " " +
                    user + ": " +
                    action +
                    '</div>');
            }
        }
    });
};