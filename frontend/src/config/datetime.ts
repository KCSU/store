import dayjs from "dayjs";
import updateLocale from "dayjs/plugin/updateLocale";
import calendar from "dayjs/plugin/calendar";
import localizedFormat from "dayjs/plugin/localizedFormat";
import "dayjs/locale/en-gb";

dayjs.extend(localizedFormat);
dayjs.extend(updateLocale);
dayjs.extend(calendar);
dayjs.locale("en-gb");
dayjs.updateLocale("en-gb", {
  calendar: {
    lastDay: "[Yesterday at] LT",
    sameDay: "[Today at] LT",
    nextDay: "[Tomorrow at] LT",
    lastWeek: "[last] dddd [at] LT",
    nextWeek: "dddd [at] LT",
    sameElse: "llll",
  },
});
