oc = oc || {};

oc.lesson = oc.lesson || {};

$(document).ready(function() {
    $('#create-lesson-form').submit(function(e) {
        e.preventDefault();
        var $lessonName = $('#lesson-name');
        var $lessonNameLbl = $('#lesson-name-lbl');
        var lessonName = $lessonName.val();
        var startingPage = $('#starting-page').val();

        $lessonName.css('background-color', 'white');
        $lessonNameLbl[0].firstChild.data = 'Lesson Name';

        var valid = false;
        if (lessonName.length == 0) {
            $lessonName.css('background-color', 'red');
            $lessonNameLbl[0].firstChild.data = 'Lesson Name - Cannot be empty.';
        } else {
            valid = true;
        }

        oc.search.page(startingPage, "#starting-page", false, function(data) {
            if (data.Result && valid) {
                $.ajax({
                    method: "POST",
                    url: "/lesson/write/",
                    data: {lesson: lessonName}
                }).done(function(data) {
                    data = JSON.parse(data);
                    if (data.Error) {
                        $lessonName.css('background-color', 'red');
                        $lessonNameLbl[0].firstChild.data = 'Lesson Name - Already exists';
                    } else {
                        window.location.href = data.Result;
                    }
                });
            }
        });
    });
});