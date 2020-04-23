var storage = new (function () {
    this.store = window.localStorage;
    this.setItem = (key, value) => {
        this.store.setItem(key, value);
    };
    this.getItem = (key) => {
        return this.store.getItem(key);
    };
    this.removeItem = (key) => {
        this.store.removeItem(key);
    };
    this.clear = () => {
        this.store.clear();
    };
});