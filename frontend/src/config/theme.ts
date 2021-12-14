import { extendTheme, theme as baseTheme } from "@chakra-ui/react";

const theme = extendTheme({
    colors: {
        brand: baseTheme.colors.purple
    },
    fonts: {
        heading: 'Montserrat',
        body: 'Montserrat'
    }
});

export default theme;