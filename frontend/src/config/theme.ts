import { extendTheme, theme as baseTheme } from "@chakra-ui/react";

const theme = extendTheme({
    colors: {
        brand: baseTheme.colors.purple,
        gray: {
            ...baseTheme.colors.gray,
            750: '#242C3A'
        }
    },
    fonts: {
        heading: 'Montserrat',
        body: 'Montserrat'
    }
});

export default theme;