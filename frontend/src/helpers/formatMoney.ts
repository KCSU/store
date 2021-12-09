export const moneyFormatter = new Intl.NumberFormat('en-GB', {
    style: 'currency',
    currency: 'GBP'
});

export function formatMoney(amount: number) {
    return moneyFormatter.format(amount);
}