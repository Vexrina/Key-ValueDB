package main;

import Interfaces.ITable;

import java.util.*;

public class Table implements ITable {
    private final Map<String, String> mapData = new HashMap<>();
    private final int time = 0;

    @Override
    public void add(String key, String value, int time) {
        this.mapData.put(key, value);
    }

    @Override
    public void delete(String key) {
        this.mapData.remove(key);
    }

    @Override
    public void put(String key, String value) {
        this.mapData.replace(key, value);
    }

    @Override
    public String get(String key) {
        return this.mapData.get(key);
    }

    @Override
    public int size(){
        return mapData.size();
    }
}
