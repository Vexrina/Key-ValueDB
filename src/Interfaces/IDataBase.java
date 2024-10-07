package Interfaces;

public interface IDataBase <K, V> {
    void add(K key, V value);
    V get(K key);
    void remove(K key);
    void put(K key, V value);
}
