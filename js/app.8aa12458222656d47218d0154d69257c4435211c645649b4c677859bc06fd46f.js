const navbarMenu = () => {
  const burger = $(".navbar-burger"),
        menu = $(".navbar-menu");

  burger.click(() => {
    [burger, menu].forEach((el) => el.toggleClass('is-active'));
  });
}

const sidebarAccordion = () => {
  const triggers = $('.trigger');

  if (triggers) {
    triggers.each((idx, el) => {
      var href = el.attributes.getNamedItem('href').nodeValue;
      href = href.slice(1);
      const collapsible = document.getElementById(href);
      new bulmaCollapsible(collapsible);
      collapsible.bulmaCollapsible('collapse');
    });
  }
}

$(() => {
  navbarMenu();
  sidebarAccordion();
});
