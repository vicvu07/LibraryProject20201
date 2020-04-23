var transparent = true;

var transparentDemo = true;
var fixedTop = false;

var navbar_initialized,
    backgroundOrange = false,
    toggle_initialized = false;

$(document).ready(function () {
    if ($(window).width() < 992) {
        nowuiKit.initRightMenu();
    }

    if ($(window).width() >= 992) {
        big_image = $('.page-header-image[data-parallax="true"]');

        $(window).on('scroll', debounce(function () {
            var current_scroll = $(this).scrollTop();

            oVal = ($(window).scrollTop() / 3);
            big_image.css({
                'transform': 'translate3d(0,' + oVal + 'px,0)',
                '-webkit-transform': 'translate3d(0,' + oVal + 'px,0)',
                '-ms-transform': 'translate3d(0,' + oVal + 'px,0)',
                '-o-transform': 'translate3d(0,' + oVal + 'px,0)'
            });

        }, 6));
    }
});

$(window).resize(function () {
    if ($(window).width() < 992) {
        nowuiKit.initRightMenu();
    }
});

nowuiKit = {
    misc: {
        navbar_menu_visible: 0
    },

    checkScrollForTransparentNavbar: debounce(function () {
        if ($(document).scrollTop() > scroll_distance) {
            if (transparent) {
                transparent = false;
                $('.navbar[color-on-scroll]').removeClass('navbar-transparent');
            }
        } else {
            if (!transparent) {
                transparent = true;
                $('.navbar[color-on-scroll]').addClass('navbar-transparent');
            }
        }
    }, 17),

    initRightMenu: function () {
        if (!toggle_initialized) {
            $toggle = $('.navbar-toggler');

            $toggle.click(function () {
                if (nowuiKit.misc.navbar_menu_visible == 1) {
                    $('html').removeClass('nav-open');
                    nowuiKit.misc.navbar_menu_visible = 0;
                    setTimeout(function () {
                        $toggle.removeClass('toggled');
                        $('#bodyClick').remove();
                    }, 550);

                } else {

                    setTimeout(function () {
                        $toggle.addClass('toggled');
                    }, 580);

                    $navbar = $(this).parent('.navbar-translate').siblings('.navbar-collapse');
                    background_image = $navbar.data('nav-image');
                    if (background_image != undefined) {
                        $navbar.css('background', "url('" + background_image + "')")
                            .removeAttr('data-nav-image')
                            .css('background-size', "cover")
                            .addClass('has-image');
                    }

                    div = '<div id="bodyClick"></div>';
                    $(div).appendTo('body').click(function () {
                        $('html').removeClass('nav-open');
                        nowuiKit.misc.navbar_menu_visible = 0;
                        setTimeout(function () {
                            $toggle.removeClass('toggled');
                            $('#bodyClick').remove();
                        }, 550);
                    });

                    $('html').addClass('nav-open');
                    nowuiKit.misc.navbar_menu_visible = 1;

                }
            });
            toggle_initialized = true;
        }
    },

    initSliders: function () {
        // Sliders for demo purpose in refine cards section
        var slider = document.getElementById('sliderRegular');

        noUiSlider.create(slider, {
            start: 40,
            connect: [true, false],
            range: {
                min: 0,
                max: 100
            }
        });

        var slider2 = document.getElementById('sliderDouble');

        noUiSlider.create(slider2, {
            start: [20, 60],
            connect: true,
            range: {
                min: 0,
                max: 100
            }
        });
    }
}

var big_image;

// Returns a function, that, as long as it continues to be invoked, will not
// be triggered. The function will be called after it stops being called for
// N milliseconds. If `immediate` is passed, trigger the function on the
// leading edge, instead of the trailing.

function debounce(func, wait, immediate) {
    var timeout;
    return function () {
        var context = this, args = arguments;
        clearTimeout(timeout);
        timeout = setTimeout(function () {
            timeout = null;
            if (!immediate) func.apply(context, args);
        }, wait);
        if (immediate && !timeout) func.apply(context, args);
    };
};
