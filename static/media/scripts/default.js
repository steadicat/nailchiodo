function getScrollXY() {
  var scrOfX = 0, scrOfY = 0;
  if( typeof( window.pageYOffset ) == 'number' ) {
    //Netscape compliant
    scrOfY = window.pageYOffset;
    scrOfX = window.pageXOffset;
  } else if( document.body && ( document.body.scrollLeft || document.body.scrollTop ) ) {
    //DOM compliant
    scrOfY = document.body.scrollTop;
    scrOfX = document.body.scrollLeft;
  } else if( document.documentElement && ( document.documentElement.scrollLeft || document.documentElement.scrollTop ) ) {
    //IE6 standards compliant mode
    scrOfY = document.documentElement.scrollTop;
    scrOfX = document.documentElement.scrollLeft;
  }
  return [ scrOfX, scrOfY ];
}

/* footnotes */

function getNote(reference) {
    return reference.parents('div').eq(0).find("#" + reference[0].id.replace("fnref", "fn"));
}
function addNoteLabel(reference) {
    var copy = reference.clone();
    copy.append('<br />');
    getNote(reference).prepend(copy);
}
function showNote(reference) {
    if (active) hideNote(active);
    active = reference;

    var note = getNote(reference);
    note.addClass("active");

    var scroll = getScrollXY();
    //note.css("left", reference.offsetLeft + "px");
    note.css("top", scroll[1]+200 + "px");
    //note.css("top", reference[0].offsetTop-50 + "px");
}
var active = null;

function hideNote(reference) {
    var note = getNote(reference);
    note.removeClass("active");
    //note.find('.footnoteReference').remove();
}

function footnotes() {
    var notes = new Array();
    $(".footnoteBackLink").remove();
    $(".footnoteReference").each(function(i) {
        this.removeAttribute("href");
        notes[i] = this;
        addNoteLabel($(this));
        $(this).hover(
            function() { showNote($(notes[i])); },
            function() { /* hideNote(notes[i]); */ }
            );
    });
}

$(footnotes);




/* languages */

function highlightOn(lang) {
    var func = function() {
        $('#' + lang).addClass('highlighted');
    };
    return func;
}
function highlightOff(lang) {
    var func = function() {
        $('#' + lang).removeClass('highlighted');
    };
    return func;
}
function unavailable() {
    var langs = arguments;
    var func = function() {
        for (var i=0; i<langs.length; i++) {
            var lang = langs[i];
            $('#' + lang).addClass('unavailable');
        }
    };
    return func;
}
function available() {
    var langs = arguments;
    var func = function() {
        for (var i=0; i<langs.length; i++) {
            var lang = langs[i];
            $('#' + lang).removeClass('unavailable');
        }
    };
    return func;
}

function languages() {
    $('.unavailable.original-en>a').hover(highlightOn('en'), highlightOff('en'));
    $('.unavailable.original-it>a').hover(highlightOn('it'), highlightOff('it'));
    $('.unavailable.original-fr>a').hover(highlightOn('fr'), highlightOff('fr'));
    $('.original-en>a').hover(unavailable('it','fr'), available('it','fr'));
    $('.original-it>a').hover(unavailable('en','fr'), available('en','fr'));
    $('.original-fr>a').hover(unavailable('en','it'), available('en','it'));
    $('.available-en>a').hover(available('en'), function() {});
    $('.available-it>a').hover(available('it'), function() {});
    $('.available-fr>a').hover(available('fr'), function() {});

    if ($('body.unavailable.original-en').length) $('#languages #en').addClass('highlighted');
    if ($('body.unavailable.original-it').length) $('#languages #it').addClass('highlighted');
    if ($('body.unavailable.original-fr').length) $('#languages #fr').addClass('highlighted');
    if ($('body.original-en').length) {
        $('#languages #it').addClass('unavailable');
        $('#languages #fr').addClass('unavailable');
    }
    if ($('body.original-it').length) {
        $('#languages #en').addClass('unavailable');
        $('#languages #fr').addClass('unavailable');
    }
    if ($('body.original-fr').length) {
        $('#languages #en').addClass('unavailable');
        $('#languages #it').addClass('unavailable');
    }
    if ($('body.available-en').length) $('#languages #en').removeClass('unavailable');
    if ($('body.available-it').length) $('#languages #it').removeClass('unavailable');
    if ($('body.available-fr').length) $('#languages #fr').removeClass('unavailable');
}

$(languages);


$(function() {
    $('.current').prepend('<span class="arrow">&#8250;</span>');
});